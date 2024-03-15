package handler

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type FilmHandler struct {
	ser *service.Service
}

func (a *FilmHandler) film(w http.ResponseWriter, req *http.Request) {
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
}

func (a *FilmHandler) getFilmsSortLike(w http.ResponseWriter, req *http.Request, orderBy, title string, desc bool) []byte {
	films, err := a.ser.Film.GetFilmsSortLike(orderBy, title, desc)
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
	films, err := a.ser.Film.GetFilmsLike(title)
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
	films, err := a.ser.Film.GetFilmsSort(orderBy, desc)
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

func (a *FilmHandler) getFilm(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		newErrorResponse(w, err, "Can't parse id from path", http.StatusBadRequest)
		return
	}

	films, err := a.ser.Film.GetFilm(id)
	if err != nil {
		newErrorResponse(w, err, "Can't get film", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		newErrorResponse(w, err, "Can't parse film to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

func (a *FilmHandler) createFilm(w http.ResponseWriter, req *http.Request) {
	isAdmin, err := a.ser.User.IsAdmin(req.Context().Value("userID").(int64))
	if err != nil {
		newErrorResponse(w, err, "", http.StatusBadRequest)
		return
	}
	if !isAdmin {
		newErrorResponse(w, errors.New("you don't have enough permissions"), "", http.StatusBadRequest)
		return
	}

	var film domain.Film
	if err := json.NewDecoder(req.Body).Decode(&film); err != nil {
		newErrorResponse(w, err, "Can't parse film from json", http.StatusBadRequest)
		return
	}
	err = a.ser.Film.CreateFilm(film)
	if err != nil {
		newErrorResponse(w, err, "Can't create film", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *FilmHandler) updateFilm(w http.ResponseWriter, req *http.Request) {
	isAdmin, err := a.ser.User.IsAdmin(req.Context().Value("userID").(int64))
	if err != nil {
		newErrorResponse(w, err, "", http.StatusBadRequest)
		return
	}
	if !isAdmin {
		newErrorResponse(w, errors.New("you don't have enough permissions"), "", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		newErrorResponse(w, err, "Can't parse id from path", http.StatusBadRequest)
		return
	}

	var film domain.Film
	if err := json.NewDecoder(req.Body).Decode(&film); err != nil {
		newErrorResponse(w, err, "Can't parse film from json", http.StatusBadRequest)
		return
	}
	film.ID = id

	err = a.ser.Film.UpdateFilm(film)
	if err != nil {
		newErrorResponse(w, err, "Can't update film", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *FilmHandler) deleteFilm(w http.ResponseWriter, req *http.Request) {
	isAdmin, err := a.ser.User.IsAdmin(req.Context().Value("userID").(int64))
	if err != nil {
		newErrorResponse(w, err, "", http.StatusBadRequest)
		return
	}
	if !isAdmin {
		newErrorResponse(w, errors.New("you don't have enough permissions"), "", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		newErrorResponse(w, err, "Can't parse id from path", http.StatusBadRequest)
		return
	}

	if err := a.ser.Film.DeleteFilm(id); err != nil {
		newErrorResponse(w, err, "Can't delete film", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *FilmHandler) getFilmActor(w http.ResponseWriter, req *http.Request, actor string) []byte {
	films, err := a.ser.Film.SearchFilmsWithActor(actor)
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
