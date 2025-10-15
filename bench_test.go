package fiberclientpool

import (
	"net/http"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

func TestMain(m *testing.M) {
	go func() {
		http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true}`))
		})
		_ = http.ListenAndServe(":8080", nil)
	}()

	time.Sleep(300 * time.Millisecond)
	os.Exit(m.Run())
}

var testURL = "http://127.0.0.1:8080/echo"

var fiberPool = NewClientPool(Config{
	Size:            runtime.GOMAXPROCS(0),
	MaxConnsPerHost: 10000,
	Timeout:         5 * time.Second,
})

var fasthttpPool = &fasthttp.LBClient{
	Clients: func() []fasthttp.BalancingClient {
		out := make([]fasthttp.BalancingClient, runtime.GOMAXPROCS(0))
		for i := range out {
			out[i] = &fasthttp.HostClient{
				Addr:         "127.0.0.1:8080",
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
				MaxConns:     10000,
			}
		}
		return out
	}(),
	Timeout: 5 * time.Second,
}

var atomicLBPool = &ClientPoolFast{
	Clients: func() []BalancingClient {
		out := make([]BalancingClient, runtime.GOMAXPROCS(0))
		for i := range out {
			out[i] = &fasthttp.HostClient{
				Addr:         "127.0.0.1:8080",
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
				MaxConns:     10000,
			}
		}
		return out
	}(),
	Timeout: 5 * time.Second,
}

func BenchmarkFiberPool_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(runtime.GOMAXPROCS(0))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := fiberPool.Next().
				R().
				SetHeader("Authorization", "Bearer 123").
				Get(testURL)
			if err != nil {
				b.Errorf("fiberpool error: %v", err)
				continue
			}
			resp.Close()
		}
	})
}

func BenchmarkFasthttpPool_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(runtime.GOMAXPROCS(0))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()

			req.Header.SetMethod("GET")
			req.SetRequestURI("/echo")
			req.Header.SetHost("127.0.0.1:8080")

			err := fasthttpPool.Do(req, resp)
			if err != nil {
				b.Errorf("fasthttp error: %v", err)
			}

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
		}
	})
}

func BenchmarkAtomicLBClient_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(runtime.GOMAXPROCS(0))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()

			req.Header.SetMethod("GET")
			req.SetRequestURI("/echo")
			req.Header.SetHost("127.0.0.1:8080")

			err := atomicLBPool.Do(req, resp)
			if err != nil {
				b.Errorf("atomic lbclient error: %v", err)
			}

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
		}
	})
}
