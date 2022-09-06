package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

func main() {

	models := []string{"../model.json", "../model_1.json"}
	sc, err := stan.Connect("test-cluster", "publisher")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()
	for i := 0; i < 10; i++ {
		//stringForByte := fmt.Sprintf("%d", i)
		//sc.Publish("chanel", []byte(stringForByte))
		//time.Sleep(100 * time.Millisecond)
	}
	for _, str := range models {
		file, err := ioutil.ReadFile(str)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = sc.Publish("chanel", file)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s %s is gone\n", "the file ", str)
		}
	}
}
