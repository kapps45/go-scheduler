package main

import (
	"go-scheduler/pkg/api"
	"go-scheduler/pkg/db"
	"go-scheduler/pkg/env"
	"log"
	"net/http"
)

const (
	envPort     = "TODO_PORT"
	defaultPort = "7540"
	defaultDB   = "scheduler.db"
	webDir      = "./web"
)

func main() {
	if err := db.Init(defaultDB); err != nil {
		log.Fatalf("ошибка базы данных: %v", err)
	}
	defer db.Close()

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := env.CheckEnv(envPort, defaultPort)
	log.Printf("сервер запущен на порту :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}
