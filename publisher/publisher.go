package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

func main() {
	models := parseDirectory()
	sc, err := stan.Connect("test-cluster", "publisher")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()
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

// Достает названия файлов из директории
func parseDirectory() (models []string) {
	path := "/Users/corkiudy/L0/models/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		models = append(models, path+f.Name())
	}
	return
}
