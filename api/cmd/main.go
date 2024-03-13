package main

import (
	handler2 "KinotekaAPI/internal/handler"
	"KinotekaAPI/internal/storage"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

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

	actorStorage := storage.NewActorStorage(conn)
	filmStorage := storage.NewFilmStorage(conn)

	handler := handler2.New(actorStorage, filmStorage)

	handler.RegisterHandlers()

	http.ListenAndServe(":8080", nil)
}
