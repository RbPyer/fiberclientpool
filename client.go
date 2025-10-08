package fiberclientpool

import (
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v3/client"
)

const (
	defaultTimeout = 60 * time.Second
	defaultSize    = 10
)

type ClientPool struct {
	pool      []*client.Client
	curNumber *atomic.Uint64
}

func newClient(cfg Config) *client.Client {
	return client.New().
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

	return &ClientPool{
		pool:      pool,
		curNumber: new(atomic.Uint64),
	}
}

func (p *ClientPool) R() *client.Client {
	return p.pool[(p.curNumber.Add(1)-1)%uint64(len(p.pool))]
}
