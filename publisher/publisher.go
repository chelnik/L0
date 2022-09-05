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
		stringForByte := fmt.Sprintf("%s %d", "boba", i)
		sc.Publish("foo", []byte(stringForByte))
		time.Sleep(1000 * time.Millisecond)
	}
}
