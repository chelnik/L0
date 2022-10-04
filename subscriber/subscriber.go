package main

import (
	model "awesomeProject"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type data struct {
	database *sql.DB
	cache    map[string]model.Model
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("../index.html")
	if err != nil {
		log.Println(err)
	}

	err = html.Execute(w, nil)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
func main() {
	database, err := openDataBase()
	checkError(err)
	// при загрузке приложения мы сразу берем данные из базы и закидываем их в мапу
	cache, err := mapFilling(database)
	if err != nil {
		log.Println(err)
	}
	dbStruct := &data{
		database: database,
		cache:    cache,
	}
	_ = dbStruct
	// создаем мультиплексор
	mux := http.NewServeMux()

	// обрабатываем пути
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/getJson", dbStruct.getJsonHandler)
	// Подписка на канал
	sc, err := stan.Connect("test-cluster", "subscriber")
	checkError(err)
	// вывод сообщения из канала
	_, err = sc.Subscribe("chanel", dbStruct.processingDataFromNats)
	checkError(err)

	// закидываем данные в кэш и в базу если их там нет
	// запуск сервера
	fmt.Println("http://localhost:4000")
	err = http.ListenAndServe("127.0.0.1:4000", mux)
	checkError(err)
}
