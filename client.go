package fiberclientpool

import (
	"github.com/gofiber/fiber/v3/client"
	"github.com/outrigdev/goid"
	"github.com/valyala/fasthttp"
)

//type cachePadded struct {
//	cursor atomic.Uint64
//	_      [56]byte
//}

type ClientPool struct {
	pool         []*client.Client
	size         int
	counters     []cachePadded
	countersSize uint64
}

func newClient(cfg Config) *client.Client {
	return client.NewWithClient(&fasthttp.Client{
		ReadTimeout:     cfg.Timeout,
		WriteTimeout:    cfg.Timeout,
		MaxConnsPerHost: int(cfg.MaxConnsPerHost),
	}).
		SetJSONMarshal(cfg.JSONMarshal).
		SetJSONUnmarshal(cfg.JSONUnmarshal).
		SetTimeout(cfg.Timeout).
		SetLogger(cfg.Logger)
}

func NewClientPool(cfg Config) *ClientPool {
	cfg.validate()
	pool := make([]*client.Client, cfg.Size)
	for i := 0; i < cfg.Size; i++ {
		pool[i] = newClient(cfg)
	}

	clientPool := &ClientPool{
		pool:         pool,
		size:         cfg.Size,
		counters:     make([]cachePadded, cfg.cursorSize),
		countersSize: uint64(cfg.cursorSize),
	}
	return clientPool
}

func (p *ClientPool) Next() *client.Client {
	n := p.counters[goid.Get()%p.countersSize].cursor.Add(1)
	return p.pool[(n-1)%uint64(p.size)]
}

func (p *ClientPool) NextWithIdx() (*client.Client, int) {
	n := p.counters[goid.Get()%p.countersSize].cursor.Add(1)
	idx := (n - 1) % uint64(p.size)
	return p.pool[idx], int(idx)
}
