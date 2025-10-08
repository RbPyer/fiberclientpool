package fiberclientpool

import (
	"encoding/json"
	stglog "log"
	"time"

	"github.com/gofiber/fiber/v3/log"
)

type Config struct {
	Size             int
	MaxConnsPerHost  int
	JSONMarshal      func(v any) ([]byte, error)
	JSONUnmarshal    func(data []byte, v any) error
	Logger           log.CommonLogger
	Timeout          time.Duration
	DisableKeepAlive bool
}

func (cfg *Config) validate() {
	if cfg.Size < 1 {
		cfg.Size = defaultSize
	}
	if cfg.JSONMarshal == nil {
		cfg.JSONMarshal = json.Marshal
	}
	if cfg.JSONUnmarshal == nil {
		cfg.JSONUnmarshal = json.Unmarshal
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}
	if cfg.Logger == nil {
		cfg.Logger = log.DefaultLogger[*stglog.Logger]()
	}
}
