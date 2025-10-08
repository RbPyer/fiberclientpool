// Package fiberclientpool provides a round-robin pool for Fiber v3 HTTP clients.
// It helps distribute requests across multiple long-lived connections,
// reducing "stickiness" to a single replica while keeping keep-alive benefits.
package fiberclientpool
