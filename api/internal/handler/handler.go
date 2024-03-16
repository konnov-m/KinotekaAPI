package handler

import (
	_ "KinotekaAPI/docs"
	"KinotekaAPI/internal/service"
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

type Handler struct {
	actor *ActorHandler
	film  *FilmHandler
	user  *UserHandler
}

func New(ser *service.Service) *Handler {
	s := &Handler{
		actor: &ActorHandler{ser: ser},
		film:  &FilmHandler{ser: ser},
		user:  &UserHandler{ser: ser},
	}

	return s
}

func (h *Handler) RegisterHandlers() {
	http.Handle("POST /actor", middlewareLog(userIdentity(http.HandlerFunc(h.actor.createActor))))
	http.Handle("GET /actor", middlewareLog(userIdentity(http.HandlerFunc(h.actor.actorsList))))

	http.Handle("GET /actor/{id}", middlewareLog(userIdentity(http.HandlerFunc(h.actor.getActor))))
	http.Handle("PUT /actor/{id}", middlewareLog(userIdentity(http.HandlerFunc(h.actor.updateActor))))
	http.Handle("DELETE /actor/{id}", middlewareLog(userIdentity(http.HandlerFunc(h.actor.deleteActor))))

	http.Handle("GET /film", middlewareLog(userIdentity(http.HandlerFunc(h.film.film))))
	http.Handle("POST /film", middlewareLog(userIdentity(http.HandlerFunc(h.film.createFilm))))

	http.Handle("GET /film/{id}", middlewareLog(userIdentity(http.HandlerFunc(h.film.getFilm))))
	http.Handle("PUT /film/{id}", middlewareLog(userIdentity(http.HandlerFunc(h.film.updateFilm))))
	http.Handle("DELETE /film/{id}", middlewareLog(userIdentity(http.HandlerFunc(h.film.deleteFilm))))
	http.Handle("POST /film/{id}", middlewareLog(userIdentity(http.HandlerFunc(h.film.addActorsToFilm))))

	http.Handle("POST /sign-up", middlewareLog(http.HandlerFunc(h.user.signUp)))
	http.Handle("POST /sign-in", middlewareLog(http.HandlerFunc(h.user.signIn)))

	http.HandleFunc("/swagger/", h.swaggerHandler)
}

func newErrorResponse(w http.ResponseWriter, err error, message string, code int) {
	var data map[string]string
	if message != "" {
		log.Printf("HTTP %d - %s. Message: %s", code, err.Error(), message)
		data = map[string]string{
			"message": message,
		}
	} else {
		log.Printf("HTTP %d - %s.", code, err.Error())
		data = map[string]string{
			"message": err.Error(),
		}
	}
	jsonData, _ := json.Marshal(data)

	http.Error(w, string(jsonData), code)
}

func (h *Handler) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	httpSwagger.WrapHandler(w, r)
}
