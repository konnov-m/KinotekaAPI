package handler

import (
	"KinotekaAPI/internal/domain"
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
}

func (a *ActorHandler) actor(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request is \"%s\". Method is %s", req.URL, req.Method)

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
	log.Printf("Request is \"%s\". Method is %s", req.URL, req.Method)

	path := req.URL.Path
	parts := strings.Split(path, "/")
	actorID := parts[len(parts)-1]
	id, err := strconv.ParseInt(actorID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
			return
		}
		jsonData, err := json.Marshal(actors)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("HTTP %d - %s", http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonData))
	} else {
		actors, err := a.s.GetActors()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
			return
		}
		jsonData, err := json.Marshal(actors)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("HTTP %d - %s", http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonData))
	}
}

func (a *ActorHandler) createActor(w http.ResponseWriter, req *http.Request) {
	var actor domain.Actor
	if err := json.NewDecoder(req.Body).Decode(&actor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}
	err := a.s.CreateActor(actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *ActorHandler) getActor(w http.ResponseWriter, req *http.Request, id int64) {
	actor, err := a.s.GetActor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}
	jsonData, err := json.Marshal(actor)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("HTTP %d - %s", http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

func (a *ActorHandler) updateActor(w http.ResponseWriter, req *http.Request, id int64) {
	var actor domain.Actor
	if err := json.NewDecoder(req.Body).Decode(&actor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}
	actor.ID = id

	err := a.s.UpdateActor(actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *ActorHandler) deleteActor(w http.ResponseWriter, req *http.Request, id int64) {
	if err := a.s.DeleteActor(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
