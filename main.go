package main

import (
	"./API"
	"./wrappers"
	"fmt"
	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"time"
)

func main() {
	var err error

	var mainAPI API.API

	loger, err := lumber.NewRotateLogger("log_"+time.Now().Format("2006-01-02")+".log", 10000, 10)
	if err != nil {
		panic("failed to setup loger")
	}

	mainAPI.Log = loger

	db, err := gorm.Open("postgres", "host=localhost port=5433 user=admin dbname=GOLibrary sslmode=disable password=12344321Qw5")
	if err != nil {
		mainAPI.Log.Error("Failed to create database connection: " + err.Error())
		panic("Failed to create database connection")
	}

	mainAPI.Db = db

	defer db.Close()

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	wrappers.HandleAll(&mainAPI)

	err = server.ListenAndServe()

	if err != nil {
		mainAPI.Log.Error("Listen server error: " + err.Error())
		fmt.Println("error:", err)
	}
}
