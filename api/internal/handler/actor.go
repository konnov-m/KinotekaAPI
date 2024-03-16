package handler

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

type ActorHandler struct {
	ser *service.Service
}

// @Summary Get List of actors
// @Security ApiKeyAuth
// @Tags actors
// @Description get list of actors
// @ID get-list-actor
// @Accept  json
// @Produce  json
// @Param withFilms query boolean false "Include films information" Enums(true,false)
// @Success 200 {object} []domain.Actor "without films"
// @Success 210 {object} []domain.ActorFilm "with films"
// @Failure 400
// @Failure 500
// @Failure default
// @Router /actor [get]
func (a *ActorHandler) actorsList(w http.ResponseWriter, req *http.Request) {
	withFilms := req.URL.Query().Get("withFilms")
	if withFilms == "true" {
		actors, err := a.ser.Actor.GetActorsWithFilms()
		if err != nil {
			newErrorResponse(w, err, "", http.StatusBadRequest)
			return
		}
		jsonData, err := json.Marshal(actors)
		if err != nil {
			newErrorResponse(w, err, "error when parse actors to json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(210)
		fmt.Fprintf(w, string(jsonData))
	} else {
		actors, err := a.ser.Actor.GetActors()
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

// @Summary Create actor
// @Security ApiKeyAuth
// @Tags actors
// @Description Create actor. You must have admin role.
// @ID create-actor
// @Accept  json
// @Produce  json
// @Param input body domain.Actor true "Actor"
// @Success 201
// @Failure 400
// @Failure default
// @Router /actor [POST]
func (a *ActorHandler) createActor(w http.ResponseWriter, req *http.Request) {
	isAdmin, err := a.ser.User.IsAdmin(req.Context().Value("userID").(int64))
	if err != nil {
		newErrorResponse(w, err, "", http.StatusBadRequest)
		return
	}
	if !isAdmin {
		newErrorResponse(w, errors.New("you don't have enough permissions"), "", http.StatusBadRequest)
		return
	}

	var actor domain.Actor
	if err := json.NewDecoder(req.Body).Decode(&actor); err != nil {
		newErrorResponse(w, err, "Can't decode actor from json", http.StatusBadRequest)
		return
	}
	err = a.ser.Actor.CreateActor(actor)
	if err != nil {
		newErrorResponse(w, err, "Can't create actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Get actor by ID
// @Security ApiKeyAuth
// @Tags actors
// @Description Get actor by ID
// @ID get-actor-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.Actor
// @Failure 400
// @Failure default
// @Router /actor/{id} [GET]
func (a *ActorHandler) getActor(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		newErrorResponse(w, err, "Can't parse id from path", http.StatusBadRequest)
		return
	}

	actor, err := a.ser.Actor.GetActor(id)
	if err != nil {
		newErrorResponse(w, err, "Can't get actor", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(actor)
	if err != nil {
		newErrorResponse(w, err, "Error when parse actor to json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

// @Summary Update actor by ID
// @Security ApiKeyAuth
// @Tags actors
// @Description Update actor by ID. You must have admin role.
// @ID update-actor-by-id
// @Accept  json
// @Produce  json
// @Param input body domain.Actor true "Actor"
// @Success 201
// @Failure 400
// @Failure default
// @Router /actor/{id} [PUT]
func (a *ActorHandler) updateActor(w http.ResponseWriter, req *http.Request) {
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

	var actor domain.Actor
	if err := json.NewDecoder(req.Body).Decode(&actor); err != nil {
		newErrorResponse(w, err, "Can't decode actor from json", http.StatusBadRequest)
		return
	}
	actor.ID = id

	err = a.ser.Actor.UpdateActor(actor)
	if err != nil {
		newErrorResponse(w, err, "Can't update actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Delete actor by ID
// @Security ApiKeyAuth
// @Tags actors
// @Description Delete actor by ID. You must have admin role.
// @ID delete-actor-by-id
// @Accept  json
// @Produce  json
// @Success 204
// @Failure 400
// @Failure default
// @Router /actor/{id} [DELETE]
func (a *ActorHandler) deleteActor(w http.ResponseWriter, req *http.Request) {
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

	if err := a.ser.Actor.DeleteActor(id); err != nil {
		newErrorResponse(w, err, "Can't delete actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
