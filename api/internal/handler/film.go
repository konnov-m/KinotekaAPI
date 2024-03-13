package handler

import (
	"KinotekaAPI/internal/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FilmHandler struct {
	s storage.FilmStorage
}

func (a *FilmHandler) film(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Query().Get("title")
	orderBy := req.URL.Query().Get("orderBy")
	desc := false
	if req.URL.Query().Get("sort") == "desc" {
		desc = true
	}
	var jsonData []byte

	if title != "" && orderBy != "" {
		films, err := a.s.GetFilmsSortLike(orderBy, title, desc)
		if err != nil {
			log.Fatal(err)
		}

		jsonData, err = json.Marshal(films)
		if err != nil {
			log.Fatal(err)
		}
	} else if title != "" {
		films, err := a.s.GetFilmsLike(title)
		if err != nil {
			log.Fatal(err)
		}

		jsonData, err = json.Marshal(films)
		if err != nil {
			log.Fatal(err)
		}
	} else if orderBy != "" {
		films, err := a.s.GetFilmsSort(orderBy, desc)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(films)
		jsonData, err = json.Marshal(films)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		films, err := a.s.GetFilms()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(films)
		jsonData, err = json.Marshal(films)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Fprintf(w, string(jsonData))
}
