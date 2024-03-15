package handler

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/service"
	"KinotekaAPI/internal/storage"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ActorHandler struct {
	s storage.ActorStorage
	u *service.UserService
}

func (a *ActorHandler) actor(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		a.actorsList(w, req)
	case http.MethodPost:
		a.createActor(w, req)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		log.Printf("Request not supported method")
	}
}

func (a *ActorHandler) actorId(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")
	actorID := parts[len(parts)-1]
	id, err := strconv.ParseInt(actorID, 10, 64)
	if err != nil {
		newErrorResponse(w, err, "Can't parse id from path", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		a.getActor(w, req, id)
	case http.MethodPut:
		a.updateActor(w, req, id)
	case http.MethodDelete:
		a.deleteActor(w, req, id)
	default:
		fmt.Fprintf(w, "Sorry, only GET, PUT and DELETE methods are supported.")
		log.Printf("Request not supported method")
	}
}

func (a *ActorHandler) actorsList(w http.ResponseWriter, req *http.Request) {
	withFilms := req.URL.Query().Get("withFilms")
	if withFilms == "true" {
		actors, err := a.s.GetActorsWithFilms()
		if err != nil {
			newErrorResponse(w, err, "Can't get actors with film", http.StatusBadRequest)
			return
		}
		jsonData, err := json.Marshal(actors)
		if err != nil {
			newErrorResponse(w, err, "Error when parse actors to json.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonData))
	} else {
		actors, err := a.s.GetActors()
		if err != nil {
			newErrorResponse(w, err, "Can't get actors", http.StatusBadRequest)
			return
		}
		jsonData, err := json.Marshal(actors)
		if err != nil {
			newErrorResponse(w, err, "Error when parse actors to json.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonData))
	}
}

func (a *ActorHandler) createActor(w http.ResponseWriter, req *http.Request) {
	var actor domain.Actor
	if err := json.NewDecoder(req.Body).Decode(&actor); err != nil {
		newErrorResponse(w, err, "Can't decode actor from json", http.StatusBadRequest)
		return
	}
	err := a.s.CreateActor(actor)
	if err != nil {
		newErrorResponse(w, err, "Can't create actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *ActorHandler) getActor(w http.ResponseWriter, req *http.Request, id int64) {
	log.Printf("getActor")
	actor, err := a.s.GetActor(id)
	if err != nil {
		newErrorResponse(w, err, "Can't get actor", http.StatusBadRequest)
		return
	}
	log.Printf("Before parse")
	jsonData, err := json.Marshal(actor)
	if err != nil {
		newErrorResponse(w, err, "Error when parse actor to json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

func (a *ActorHandler) updateActor(w http.ResponseWriter, req *http.Request, id int64) {
	var actor domain.Actor
	if err := json.NewDecoder(req.Body).Decode(&actor); err != nil {
		newErrorResponse(w, err, "Can't decode actor from json", http.StatusBadRequest)
		return
	}
	actor.ID = id

	err := a.s.UpdateActor(actor)
	if err != nil {
		newErrorResponse(w, err, "Can't update actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *ActorHandler) deleteActor(w http.ResponseWriter, req *http.Request, id int64) {
	if err := a.s.DeleteActor(id); err != nil {
		newErrorResponse(w, err, "Can't delete actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
