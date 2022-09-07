package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type Post struct {
	ID       int
	Text     string
	Author   string
	Comments int
	Time     time.Time
}

func handle(w http.ResponseWriter, req *http.Request) {
	s := ""
	for i := 0; i < 1000; i++ {
		p := &Post{ID: i, Text: "yo post", Time: time.Now()}
		s += fmt.Sprintf("%#v", p)
	}
	w.Write([]byte(s))
}
func main() {
	http.HandleFunc("/", handle)
	fmt.Println("http://localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
// 21 page
// ab -t 300 -n 1000000000 -c 10 http://127.0.0.1:8080/
// ab -t 300 -n 100000 -c 10 http://127.0.0.1:8080