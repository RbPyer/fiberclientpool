package fiberclientpool

import (
	"encoding/json"
	stglog "log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v3/log"
)

const defaultTimeout = 60 * time.Second

type Config struct {
	Size             int
	CursorSize       int
	MaxConnsPerHost  int64
	JSONMarshal      func(v any) ([]byte, error)
	JSONUnmarshal    func(data []byte, v any) error
	Logger           log.CommonLogger
	Timeout          time.Duration
	DisableKeepAlive bool
}

func (cfg *Config) validate() {
	if cfg.Size < 1 {
		cfg.Size = runtime.GOMAXPROCS(0)
	}
	if cfg.CursorSize < 1 {
		cfg.CursorSize = cfg.Size
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
