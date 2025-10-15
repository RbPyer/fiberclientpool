package fiberclientpool

import "testing"

var testCfg = Config{}

func BenchmarkClientPool_Next_Parallel(b *testing.B) {
	pool := NewClientPool(testCfg)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.Next()
		}
	})
}
