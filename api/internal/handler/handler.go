package handler

import (
	"KinotekaAPI/internal/storage"
	"net/http"
)

type Handler struct {
	actor *ActorHandler
	film  *FilmHandler
}

func New(a storage.ActorStorage, f storage.FilmStorage) *Handler {
	s := &Handler{
		actor: &ActorHandler{s: a},
		film:  &FilmHandler{s: f},
	}

	return s
}

func (h *Handler) RegisterHandlers() {
	http.HandleFunc("/actor", h.actor.actorsList)
	http.HandleFunc("/film", h.film.film)
}
