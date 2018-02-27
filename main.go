package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/XpertVision/Library/api"
	"github.com/XpertVision/Library/wrappers"
)

func main() {
	var err error

	var mainAPI api.API

	loger, err := lumber.NewRotateLogger("log_"+time.Now().Format("2006-01-02")+".log", 10000, 10)
	if err != nil {
		panic("failed to setup loger")
	}

	db, err := gorm.Open("postgres", "host=localhost port=5433 user=admin dbname=GOLibrary sslmode=disable password=12344321Qw5")
	if err != nil {
		mainAPI.Log.Error("Failed to create database connection: " + err.Error())
		panic("Failed to create database connection")
	}

	api.New(db, loger)

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
