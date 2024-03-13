package handler

import (
	"KinotekaAPI/internal/storage"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type ActorHandler struct {
	s storage.ActorStorage
}

func (a *ActorHandler) actorsList(w http.ResponseWriter, req *http.Request) {
	actors, err := a.s.GetActors()
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := json.Marshal(actors)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(jsonData))
}
