package service

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/storage"
)

type actorService struct {
	s storage.ActorStorage
}

func NewActorService(s storage.ActorStorage) Actor {
	return &actorService{
		s: s,
	}
}

func (a *actorService) GetActorsWithFilms() ([]domain.ActorFilm, error) {
	return a.s.GetActorsWithFilms()
}

func (a *actorService) GetActors() ([]domain.Actor, error) {
	return a.s.GetActors()
}

func (a *actorService) CreateActor(actor domain.Actor) error {
	return a.s.CreateActor(actor)
}

func (a *actorService) GetActor(id int64) (domain.Actor, error) {
	return a.s.GetActor(id)
}

func (a *actorService) UpdateActor(actor domain.Actor) error {
	return a.s.UpdateActor(actor)
}

func (a *actorService) DeleteActor(id int64) error {
	if err := a.s.DeleteActorsFilms(id); err != nil {
		return err
	}

	return a.s.DeleteActor(id)
}
