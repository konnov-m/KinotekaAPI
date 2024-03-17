package service

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/storage"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type User interface {
	CreateUser(user domain.User, role string) error
	GenerateToken(login, password string) (string, error)
	IsAdmin(id int64) (bool, error)
}

type Actor interface {
	GetActorsWithFilms() ([]domain.ActorFilm, error)
	CreateActor(actor domain.Actor) error
	GetActor(id int64) (domain.Actor, error)
	UpdateActor(actor domain.Actor) error
	DeleteActor(id int64) error
	GetActors() ([]domain.Actor, error)
}

type Film interface {
	GetFilmsSortLike(orderBy, title string, desc bool) ([]domain.Film, error)
	GetFilmsLike(title string) ([]domain.Film, error)
	GetFilmsSort(orderBy string, desc bool) ([]domain.Film, error)
	GetFilm(id int64) (domain.Film, error)
	CreateFilm(a domain.Film) error
	UpdateFilm(a domain.Film) error
	DeleteFilm(id int64) error
	SearchFilmsWithActor(substr string) ([]domain.ActorFilm, error)
	AddActorToFilm(filmId int64, actorId []int64) error
}

type Service struct {
	User
	Actor
	Film
}

func NewService(s *storage.Storage) *Service {
	return &Service{
		User:  NewUserService(s.UserStorage),
		Actor: NewActorService(s.ActorStorage),
		Film:  NewFilmService(s.FilmStorage),
	}
}
