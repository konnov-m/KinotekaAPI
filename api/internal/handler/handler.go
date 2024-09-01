package handler

import (
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "kinoteka/docs"
	"kinoteka/internal/service"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	actor *ActorHandler
	film  *FilmHandler
	user  *UserHandler
	ser   *service.Service
}

func New(ser *service.Service) *Handler {
	s := &Handler{
		actor: &ActorHandler{ser: ser},
		film:  &FilmHandler{ser: ser},
		user:  &UserHandler{ser: ser},
		ser:   ser,
	}

	return s
}

func (h *Handler) RegisterHandlers() {
	http.Handle("POST /actor", middlewareLog(h.userIdentity(http.HandlerFunc(h.actor.createActor))))
	http.Handle("GET /actor", middlewareLog(h.userIdentity(http.HandlerFunc(h.actor.actorsList))))

	http.Handle("GET /actor/{id}", middlewareLog(h.userIdentity(http.HandlerFunc(h.actor.getActor))))
	http.Handle("PUT /actor/{id}", middlewareLog(h.userIdentity(http.HandlerFunc(h.actor.updateActor))))
	http.Handle("DELETE /actor/{id}", middlewareLog(h.userIdentity(http.HandlerFunc(h.actor.deleteActor))))

	http.Handle("GET /film", middlewareLog(h.userIdentity(http.HandlerFunc(h.film.film))))
	http.Handle("POST /film", middlewareLog(h.userIdentity(http.HandlerFunc(h.film.createFilm))))

	http.Handle("GET /film/{id}", middlewareLog(h.userIdentity(http.HandlerFunc(h.film.getFilm))))
	http.Handle("PUT /film/{id}", middlewareLog(h.userIdentity(http.HandlerFunc(h.film.updateFilm))))
	http.Handle("DELETE /film/{id}", middlewareLog(h.userIdentity(http.HandlerFunc(h.film.deleteFilm))))
	http.Handle("POST /film/{id}", middlewareLog(h.userIdentity(http.HandlerFunc(h.film.addActorsToFilm))))

	http.Handle("POST /sign-up", middlewareLog(http.HandlerFunc(h.user.signUp)))
	http.Handle("POST /sign-in", middlewareLog(http.HandlerFunc(h.user.signIn)))

	http.Handle("/metrics", promhttp.Handler())

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
