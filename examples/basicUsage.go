package main

import (
	"fmt"
	"time"

	pool "github.com/RbPyer/fiberclientpool"
)

func main() {
	p := pool.NewClientPool(pool.Config{
		Size:             10,
		Timeout:          10 * time.Second,
		MaxConnsPerHost:  1,
		DisableKeepAlive: false,
	})

	resp, err := p.R().
		Debug().
		SetHeader("Authorization", "Bearer 123").
		Get("https://reqbin.com/echo")
	if err != nil {
		panic(err)
	}
	fmt.Println("Status:", resp.StatusCode())
}
