package storage

import (
	"KinotekaAPI/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ActorStorage interface {
	GetActors() ([]domain.Actor, error)
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
