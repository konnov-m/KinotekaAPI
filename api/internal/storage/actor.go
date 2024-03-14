package storage

import (
	"KinotekaAPI/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ActorStorage interface {
	GetActors() ([]domain.Actor, error)
	CreateActor(a domain.Actor) error
	GetActor(id int64) (domain.Actor, error)
}

type actorStorage struct {
	db *sqlx.DB
}

func NewActorStorage(conn *sqlx.DB) ActorStorage {
	return &actorStorage{
		db: conn,
	}
}

const getActors = `SELECT * FROM actors`

func (s *actorStorage) GetActors() ([]domain.Actor, error) {
	var actors []domain.Actor
	err := s.db.Select(&actors, getActors)
	if err != nil {
		return nil, err
	}

	return actors, nil
}

const getActor = `SELECT * FROM actors WHERE id = $1`

func (s *actorStorage) GetActor(id int64) (domain.Actor, error) {
	var actor domain.Actor
	err := s.db.Get(&actor, getActor, id)
	if err != nil {
		return domain.Actor{}, err
	}

	return actor, nil
}

const saveActor = `INSERT INTO actors (name, surname, patronymic, birthday, sex, information)
VALUES ($1, $2, $3, $4, $5, $6);`

func (s *actorStorage) CreateActor(a domain.Actor) error {
	_, err := s.db.Exec(saveActor, a.Name, a.Surname, a.Patronymic, a.Birthday, a.Sex, a.Information)
	if err != nil {
		return err
	}

	return nil
}
