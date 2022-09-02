package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"time"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher")
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 100; i++ {
		sc.Publish("foo", []byte("hello vadim"))
		time.Sleep(1000 * time.Millisecond)
	}
}
