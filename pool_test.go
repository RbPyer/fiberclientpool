package fiberclientpool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	fiberclient "github.com/gofiber/fiber/v3/client"
)

func newTestPool(size int, timeout time.Duration) *ClientPool {
	cfg := Config{
		Size:    size,
		Timeout: timeout,
	}
	return NewClientPool(cfg)
}

func TestClientPool_RoundRobinOrder(t *testing.T) {
	t.Parallel()

	size := 4
	p := newTestPool(size, 3*time.Second)

	seen := make([]*fiberclient.Client, 0, size*2)
	for i := 0; i < size*2; i++ {
		seen = append(seen, p.Next())
	}

	uniq := map[*fiberclient.Client]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := uniq[seen[i]]; ok {
			t.Fatalf("duplicate client in first cycle at i=%d", i)
		}
		uniq[seen[i]] = struct{}{}
	}

	for i := 0; i < size; i++ {
		if seen[i] != seen[size+i] {
			t.Fatalf("round-robin mismatch at pos %d: %p != %p", i, seen[i], seen[size+i])
		}
	}
}

func TestClientPool_NextWithIdx(t *testing.T) {
	t.Parallel()

	size := 5
	p := newTestPool(size, 2*time.Second)

	for i := 0; i < size*3; i++ {
		c, idx := p.NextWithIdx()
		if c == nil {
			t.Fatalf("got nil client at iteration %d", i)
		}
		if idx < 0 || idx >= size {
			t.Fatalf("invalid idx=%d (out of range)", idx)
		}
	}
}

func TestClientPool_ConcurrentAccessIsSafe(t *testing.T) {
	t.Parallel()

	p := newTestPool(8, 3*time.Second)

	const goroutines = 64
	const iters = 10_000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	var nilCount atomic.Int64
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				c := p.Next()
				if c == nil {
					nilCount.Add(1)
				}
			}
		}()
	}
	wg.Wait()

	if got := nilCount.Load(); got != 0 {
		t.Fatalf("Next() returned nil %d times", got)
	}
}
