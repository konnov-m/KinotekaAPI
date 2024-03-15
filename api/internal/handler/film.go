package handler

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/service"
	"KinotekaAPI/internal/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type FilmHandler struct {
	s storage.FilmStorage
	u *service.UserService
}

func (a *FilmHandler) film(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		title := req.URL.Query().Get("title")
		orderBy := req.URL.Query().Get("orderBy")
		actor := req.URL.Query().Get("actor")
		desc := false
		if req.URL.Query().Get("sort") == "desc" {
			desc = true
		}
		var jsonData []byte

		if title != "" && orderBy != "" {
			jsonData = a.getFilmsSortLike(w, req, orderBy, title, desc)
		} else if title != "" {
			jsonData = a.getFilmsLike(w, req, title)
		} else if actor != "" {
			jsonData = a.getFilmActor(w, req, actor)
		} else {
			jsonData = a.getFilmsSort(w, req, orderBy, desc)
		}

		if jsonData != nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, string(jsonData))
		}
	case http.MethodPost:
		a.createFilm(w, req)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		log.Printf("Request not supported method")
	}
}

func (a *FilmHandler) filmId(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")
	filmID := parts[len(parts)-1]
	id, err := strconv.ParseInt(filmID, 10, 64)
	if err != nil {
		newErrorResponse(w, err, "Can't parse id from path", http.StatusBadRequest)
		return
	}
	var jsonData []byte

	switch req.Method {
	case http.MethodGet:
		jsonData = a.getFilm(w, req, id)
	case http.MethodPut:
		a.updateFilm(w, req, id)
	case http.MethodDelete:
		a.deleteFilm(w, req, id)
	default:
		fmt.Fprintf(w, "Sorry, only GET, PUT and DELETE methods are supported.")
		log.Printf("Request not supported method")
	}
	if jsonData != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonData))
	}
}

func (a *FilmHandler) getFilmsSortLike(w http.ResponseWriter, req *http.Request, orderBy, title string, desc bool) []byte {
	films, err := a.s.GetFilmsSortLike(orderBy, title, desc)
	if err != nil {
		newErrorResponse(w, err, "Can't get sort films", http.StatusBadRequest)
		return nil
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		newErrorResponse(w, err, "Can't parse films to json", http.StatusInternalServerError)
		return nil
	}
	return jsonData
}

func (a *FilmHandler) getFilmsLike(w http.ResponseWriter, req *http.Request, title string) []byte {
	films, err := a.s.GetFilmsLike(title)
	if err != nil {
		newErrorResponse(w, err, "Can't get films", http.StatusBadRequest)
		return nil
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		newErrorResponse(w, err, "Can't parse films to json", http.StatusInternalServerError)
		return nil
	}

	return jsonData
}

func (a *FilmHandler) getFilmsSort(w http.ResponseWriter, req *http.Request, orderBy string, desc bool) []byte {
	films, err := a.s.GetFilmsSort(orderBy, desc)
	if err != nil {
		newErrorResponse(w, err, "Can't get films", http.StatusBadRequest)
		return nil
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		newErrorResponse(w, err, "Can't parse films to json", http.StatusInternalServerError)
		return nil
	}

	return jsonData
}

func (a *FilmHandler) getFilm(w http.ResponseWriter, req *http.Request, id int64) []byte {
	films, err := a.s.GetFilm(id)
	if err != nil {
		newErrorResponse(w, err, "Can't get film", http.StatusBadRequest)
		return nil
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		newErrorResponse(w, err, "Can't parse film to json", http.StatusInternalServerError)
		return nil
	}

	return jsonData
}

func (a *FilmHandler) createFilm(w http.ResponseWriter, req *http.Request) {
	var film domain.Film
	if err := json.NewDecoder(req.Body).Decode(&film); err != nil {
		newErrorResponse(w, err, "Can't parse film from json", http.StatusBadRequest)
		return
	}
	err := a.s.CreateFilm(film)
	if err != nil {
		newErrorResponse(w, err, "Can't create film", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *FilmHandler) updateFilm(w http.ResponseWriter, req *http.Request, id int64) {
	var film domain.Film
	if err := json.NewDecoder(req.Body).Decode(&film); err != nil {
		newErrorResponse(w, err, "Can't parse film from json", http.StatusBadRequest)
		return
	}
	film.ID = id

	err := a.s.UpdateFilm(film)
	if err != nil {
		newErrorResponse(w, err, "Can't update film", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *FilmHandler) deleteFilm(w http.ResponseWriter, req *http.Request, id int64) {
	if err := a.s.DeleteFilm(id); err != nil {
		newErrorResponse(w, err, "Can't delete film", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *FilmHandler) getFilmActor(w http.ResponseWriter, req *http.Request, actor string) []byte {
	films, err := a.s.SearchFilmsWithActor(actor)
	if err != nil {
		newErrorResponse(w, err, "Can't get films with actor", http.StatusBadRequest)
		return nil
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		newErrorResponse(w, err, "Can't parse film to json", http.StatusInternalServerError)
		return nil
	}

	return jsonData
}
