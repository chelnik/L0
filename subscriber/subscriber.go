package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"net/http"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("hello po"))
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	sc, err := stan.Connect("test-cluster", "consumer")
	if err != nil {
		fmt.Println(err)
	}
	_, err = sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		fmt.Println(err)
	}
	err = http.ListenAndServe("127.0.0.1:4000", mux)
	if err != nil {
		fmt.Println(err)
	}

}

// consumer is subscriber
