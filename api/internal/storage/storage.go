package storage

import (
	"KinotekaAPI/internal/domain"
	"github.com/jmoiron/sqlx"
)

type FilmStorage interface {
	GetFilmsLike(title string) ([]domain.Film, error)
	GetFilmsSort(orderBy string, desc bool) ([]domain.Film, error)
	GetFilmsSortLike(orderBy, title string, desc bool) ([]domain.Film, error)
	GetFilm(id int64) (domain.Film, error)
	CreateFilm(a domain.Film) error
	UpdateFilm(a domain.Film) error
	DeleteFilm(id int64) error
	SearchFilmsWithActor(substr string) ([]domain.ActorFilm, error)
	DeleteFilmsActors(id int64) error
}

type ActorStorage interface {
	GetActors() ([]domain.Actor, error)
	CreateActor(a domain.Actor) error
	GetActor(id int64) (domain.Actor, error)
	UpdateActor(a domain.Actor) error
	DeleteActor(id int64) error
	GetActorsWithFilms() ([]domain.ActorFilm, error)
	DeleteActorsFilms(id int64) error
}

type UserStorage interface {
	GetUser(login, pass string) (*domain.User, error)
	CreateUser(user domain.User, role string) error
	GetRole(userId int64) ([]domain.Role, error)
}

type Storage struct {
	FilmStorage
	ActorStorage
	UserStorage
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		FilmStorage:  NewFilmStorage(db),
		ActorStorage: NewActorStorage(db),
		UserStorage:  NewUserStorage(db),
	}
}
