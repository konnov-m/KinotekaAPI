package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"kinoteka/internal/domain"
	"kinoteka/internal/service"
	"net/http"
	"strconv"
)

type FilmHandler struct {
	ser *service.Service
}

// @Summary Get List of films
// @Security ApiKeyAuth
// @Tags films
// @Description get list of films. You can search by title, actor name, surname
// @ID get-list-films
// @Accept  json
// @Produce  json
// @Param title query string false "Search by title"
// @Param sort query string false "Sort list by desc or asc" Enums(desc,asc)
// @Param orderBy query string false "sort by params" Enums(rating,title,year)
// @Param actor query string false "Search by actor"
// @Success 200 {object} []domain.Film
// @Success 210 {object} []domain.ActorFilm
// @Failure 400
// @Failure 500
// @Failure default
// @Router /film [get]
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(210)
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

// @Summary Get Film by ID
// @Security ApiKeyAuth
// @Tags films
// @Description Get Film by ID
// @ID get-film-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.Actor
// @Failure 400
// @Failure default
// @Router /film/{id} [GET]
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

// @Summary Create Film
// @Security ApiKeyAuth
// @Tags films
// @Description Create Film. You must have admin role.
// @ID create-film
// @Accept  json
// @Produce  json
// @Param input body domain.Film true "Film"
// @Success 201
// @Failure 400
// @Failure 500
// @Failure default
// @Router /film [POST]
func (a *FilmHandler) createFilm(w http.ResponseWriter, req *http.Request) {
	isAdmin, err := a.ser.User.IsAdmin(req.Context().Value("userID").(int64))
	if err != nil {
		newErrorResponse(w, err, "Not admin", http.StatusBadRequest)
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

// @Summary Update Film by ID
// @Security ApiKeyAuth
// @Tags films
// @Description Update Film by ID. You must have admin role.
// @ID update-film-by-id
// @Accept  json
// @Produce  json
// @Param input body domain.Film true "Film"
// @Success 201
// @Failure 400
// @Failure default
// @Router /film/{id} [PUT]
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

// @Summary Delete Film by ID
// @Security ApiKeyAuth
// @Tags films
// @Description Delete Film by ID. You must have admin role.
// @ID delete-film-by-id
// @Accept  json
// @Produce  json
// @Param input body domain.Film true "Film"
// @Success 201
// @Failure 400
// @Failure default
// @Router /film/{id} [DELETE]
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

type Data struct {
	Actors []int64 `json:"actors"`
}

// @Summary Add actors to film by id
// @Security ApiKeyAuth
// @Tags films
// @Description Add actors to film by id. You must have admin role.
// @ID add-actor-to-film-by-id
// @Accept  json
// @Produce  json
// @Param input body Data true "Array of actor's id"
// @Success 201
// @Failure 400
// @Failure default
// @Router /film/{id} [POST]
func (a *FilmHandler) addActorsToFilm(w http.ResponseWriter, req *http.Request) {
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

	var data Data
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		newErrorResponse(w, err, "Can't parse data from json", http.StatusBadRequest)
		return
	}
	if data.Actors == nil {
		newErrorResponse(w, errors.New("data.Actors is nil"), "Wrong input form", http.StatusBadRequest)
		return
	}

	err = a.ser.Film.AddActorToFilm(id, data.Actors)
	if err != nil {
		newErrorResponse(w, err, "Can't add actor to film", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
