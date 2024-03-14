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
	log.Printf("Request is \"%s\". Method is %s", req.URL.Path, req.Method)

	title := req.URL.Query().Get("title")
	orderBy := req.URL.Query().Get("orderBy")
	desc := false
	if req.URL.Query().Get("sort") == "desc" {
		desc = true
	}
	var jsonData []byte

	if title != "" && orderBy != "" {
		jsonData = a.getFilmsSortLike(w, req, orderBy, title, desc)
	} else if title != "" {
		jsonData = a.getFilmsLike(w, req, title)
	} else {
		jsonData = a.getFilmsSort(w, req, orderBy, desc)
	}

	if jsonData != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonData))
	}
}

func (a *FilmHandler) getFilmsSortLike(w http.ResponseWriter, req *http.Request, orderBy, title string, desc bool) []byte {
	films, err := a.s.GetFilmsSortLike(orderBy, title, desc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return nil
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("HTTP %d - %s", http.StatusInternalServerError, err.Error())
		return nil
	}
	return jsonData
}

func (a *FilmHandler) getFilmsLike(w http.ResponseWriter, req *http.Request, title string) []byte {
	films, err := a.s.GetFilmsLike(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return nil
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("HTTP %d - %s", http.StatusInternalServerError, err.Error())
		return nil
	}

	return jsonData
}

func (a *FilmHandler) getFilmsSort(w http.ResponseWriter, req *http.Request, orderBy string, desc bool) []byte {
	films, err := a.s.GetFilmsSort(orderBy, desc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return nil
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("HTTP %d - %s", http.StatusInternalServerError, err.Error())
		return nil
	}

	return jsonData
}
