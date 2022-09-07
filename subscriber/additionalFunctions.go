package main

import (
	model "awesomeProject"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
)

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

// Записывает данные в кэш и базу данных из данных которые отдал publisher
func (dbStruct *data) processingDataFromNats(message *stan.Msg) {
	var temp model.Model
	err := json.Unmarshal(message.Data, &temp)
	checkError(err)

	if temp.Order_uid == "" {
		fmt.Println("пустой uid")
		return
	}
	if _, ok := dbStruct.cache[temp.Order_uid]; ok {
		fmt.Println("данный uid уже существует", dbStruct.cache[temp.Order_uid].Order_uid)
		return
	}
	dbStruct.cache[temp.Order_uid] = temp
	_, err = dbStruct.database.Exec("INSERT INTO orders (Order_uid, DataJson) VALUES ($1, $2)", temp.Order_uid, message.Data)
	if err != nil {
		log.Println(err, "error nats")
	}
	fmt.Printf("Добавил %s в кэш и базу из publisher\n", temp.Order_uid)
}

// Обработка получения json по ключу
func (dbStruct *data) getJsonHandler(w http.ResponseWriter, r *http.Request) {
	signature := r.FormValue("signature")
	if dbStruct.cache[signature].Order_uid == signature && signature != "" {
		html, err := template.ParseFiles("../viewJson.html")
		if err != nil {
			fmt.Println(err)
		}
		err = html.Execute(w, dbStruct.cache[signature])
	} else {
		// обработка невалидного ключа
		_, err := w.Write([]byte("NOT FOUND\nPlease enter a valid key"))
		if err != nil {
			fmt.Println(err)
		}
	}
}
