package main

import (
	model "awesomeProject"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
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

// здесь сделаем еще один обработчик
// по нажатию на кнопку пользователя будет перекидывать на другую страницу и выводить либо ошибку либо значение json
// соответственно будем получать значение поля сабмит и если мы получили данные из кэша выводим их

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
	// обрабатываем путь
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/getJson", func(w http.ResponseWriter, r *http.Request) {
		signature := r.FormValue("signature")
		if dbStruct.cache[signature].Order_uid == signature && signature != "" {
			html, err := template.ParseFiles("../viewJson.html")
			if err != nil {
				fmt.Println(err)
			}
			err = html.Execute(w, dbStruct.cache[signature])
			fmt.Println(dbStruct.cache[signature].Order_uid, dbStruct.cache[signature].Order_uid)
		} else {
			_, err = w.Write([]byte("NOT FOUND\nPlease enter a valid key"))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(dbStruct.cache[signature].Order_uid, dbStruct.cache[signature].Order_uid)
		}

	})
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

func (dbStruct *data) processingDataFromNats(message *stan.Msg) {
	var temp model.Model
	err := json.Unmarshal(message.Data, &temp)
	checkError(err)
	if temp.Order_uid == "" {
		fmt.Println("пустой uid")
		return
	}
	if dbStruct.cache[temp.Order_uid].Order_uid != "" {
		fmt.Println("данный uid уже существует")
		fmt.Println(dbStruct.cache[temp.Order_uid].Order_uid)
		return
	}

	dbStruct.cache[temp.Order_uid] = temp
	_, err = dbStruct.database.Exec("INSERT INTO orders (Order_uid, DataJson)VALUES ($1, $2)", temp.Order_uid, message.Data)
	fmt.Printf("Добавил %s в кэш и базу из publisher\n", temp.Order_uid)
	if err != nil {
		log.Println(err, "error nats")
	}
}

// Открывает базу данных
func openDataBase() (*sql.DB, error) {
	database, err := sql.Open("postgres", "postgresql://corkiudy:test@127.0.0.1:5432/wildberries?sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err = database.Ping(); err != nil { // проверка того что все настроено правильно
		return nil, err
	}
	return database, nil
}

// Заполняет кэш из базы данных
func mapFilling(database *sql.DB) (map[string]model.Model, error) {
	var key string
	var value model.Model
	var sliceForJson []byte
	var mapa map[string]model.Model = make(map[string]model.Model)
	rows, err := database.Query("SELECT * FROM orders")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s scan", err))
	}
	for rows.Next() {
		err = rows.Scan(&key, &sliceForJson)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s mapFilling", err))
		}
		json.Unmarshal(sliceForJson, &value)
		mapa[key] = value
		fmt.Printf("Добавил %s в кэш из базы\n", key)
	}
	return mapa, nil
}
