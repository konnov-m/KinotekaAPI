package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	handler2 "kinoteka/internal/handler"
	"kinoteka/internal/service"
	"kinoteka/internal/storage"
	"log"
	"net/http"
	"os"
)

// @title Kinoteka API
// @version 1.0
// @description API Server for kinoteka Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	name := os.Getenv("PG_NAME")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	host := os.Getenv("PG_HOST")

	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, pass, host, name)
	//url := "postgres://admin:1234567890@db:5432/kinoteka_api?sslmode=disable"

	conn, err := sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	storages := storage.NewStorage(conn)
	services := service.NewService(storages)

	handler := handler2.New(services)

	handler.RegisterHandlers()

	http.ListenAndServe(":8080", nil)
}
