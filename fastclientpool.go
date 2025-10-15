package fiberclientpool

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	_ "unsafe"

	//"github.com/outrigdev/goid"
	"github.com/valyala/fasthttp"
)

//go:linkname procPin runtime.procPin
func procPin() int

//go:linkname procUnpin runtime.procUnpin
func procUnpin()

type BalancingClient interface {
	DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error
}

type cachePadded struct {
	cursor atomic.Uint64
	_      [56]byte
}

type ClientPoolFast struct {
	Clients  []BalancingClient
	cs       []*lbClient
	counters []cachePadded
	Timeout  time.Duration
	size     uint64

	countersSize uint64

	once sync.Once
}

const DefaultLBClientTimeout = time.Second

func (cc *ClientPoolFast) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	return cc.get().DoDeadline(req, resp, deadline)
}

func (cc *ClientPoolFast) DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	return cc.get().DoDeadline(req, resp, deadline)
}

func (cc *ClientPoolFast) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	timeout := cc.Timeout
	if timeout <= 0 {
		timeout = DefaultLBClientTimeout
	}
	return cc.DoTimeout(req, resp, timeout)
}

func (cc *ClientPoolFast) init() {
	if len(cc.Clients) == 0 {
		panic("BUG: ClientPoolFast.Clients cannot be empty")
	}

	cc.cs = make([]*lbClient, len(cc.Clients))
	for i, c := range cc.Clients {
		cc.cs[i] = &lbClient{c: c}
	}

	ncpu := runtime.GOMAXPROCS(0)
	if ncpu <= 0 {
		ncpu = 1
	}
	cc.counters = make([]cachePadded, ncpu)
	cc.countersSize = uint64(ncpu)
	cc.size = uint64(len(cc.Clients))
}

func (cc *ClientPoolFast) get() *lbClient {
	cc.once.Do(cc.init)

	pid := uint64(procPin()) % cc.countersSize
	n := cc.counters[pid].cursor.Add(1)
	procUnpin()

	idx := (n - 1) % cc.size
	return cc.cs[idx]
}

type lbClient struct {
	c BalancingClient
}

func (c *lbClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	return c.c.DoDeadline(req, resp, deadline)
}
