// Package fiberclientpool provides a high-performance pool for Fiber v3 HTTP clients.
//
// The pool distributes requests across multiple long-lived *fiber.Client*
// instances to balance load and preserve connection reuse (keep-alive).
//
// Internally, it uses a sharded counter approach: each logical processor (P)
// has its own atomic counter, reducing cache-line contention and avoiding
// global atomic bottlenecks under high concurrency.
//
// This design favors throughput and low contention over strict round-robin
// ordering. In practice, requests are evenly distributed, but when the number
// of goroutines greatly exceeds GOMAXPROCS, multiple goroutines may share the
// same counter-shard, resulting in approximate rather than perfectly sequential
// round-robin behavior.
//
// The pool is safe for concurrent use by multiple goroutines.
package fiberclientpool
