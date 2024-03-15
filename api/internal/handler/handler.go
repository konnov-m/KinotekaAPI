package handler

import (
	"KinotekaAPI/internal/service"
	"KinotekaAPI/internal/storage"
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	actor *ActorHandler
	film  *FilmHandler
	user  *UserHandler
}

func New(a storage.ActorStorage, f storage.FilmStorage, u *service.UserService) *Handler {
	s := &Handler{
		actor: &ActorHandler{s: a, u: u},
		film:  &FilmHandler{s: f, u: u},
		user:  &UserHandler{s: u},
	}

	return s
}

func (h *Handler) RegisterHandlers() {
	http.Handle("/actor", middlewareLog(userIdentity(http.HandlerFunc(h.actor.actor))))
	http.Handle("/actor/", middlewareLog(userIdentity(http.HandlerFunc(h.actor.actorId))))

	http.Handle("/film", middlewareLog(userIdentity(http.HandlerFunc(h.film.film))))
	http.Handle("/film/", middlewareLog(userIdentity(http.HandlerFunc(h.film.filmId))))

	http.Handle("/sign-up", middlewareLog(http.HandlerFunc(h.user.signUp)))
	http.Handle("/sign-in", middlewareLog(http.HandlerFunc(h.user.signIn)))
}

func newErrorResponse(w http.ResponseWriter, err error, message string, code int) {
	data := map[string]string{
		"message": message,
	}
	jsonData, _ := json.Marshal(data)

	http.Error(w, string(jsonData), code)
	log.Printf("HTTP %d - %s. Message: %s", code, err.Error(), message)
}
