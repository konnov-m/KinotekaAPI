package storage

import (
	"KinotekaAPI/internal/domain"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type ActorStorage interface {
	GetActors() ([]domain.Actor, error)
	CreateActor(a domain.Actor) error
	GetActor(id int64) (domain.Actor, error)
	UpdateActor(a domain.Actor) error
	DeleteActor(id int64) error
	GetActorsWithFilms() ([]domain.ActorFilm, error)
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

const updateActor = `UPDATE actors SET name=$1, surname=$2, patronymic=$3, birthday=$4, sex=$5, information=$6
WHERE id=$7;`

func (s *actorStorage) UpdateActor(a domain.Actor) error {
	_, err := s.db.Exec(updateActor, a.Name, a.Surname, a.Patronymic, a.Birthday, a.Sex, a.Information, a.ID)
	if err != nil {
		return err
	}

	return nil
}

const deleteActor = `DELETE from actors WHERE id=$1;`

func (s *actorStorage) DeleteActor(id int64) error {
	_, err := s.db.Exec(deleteActor, id)
	if err != nil {
		return err
	}

	return nil
}

const getActorsWithFilms = `SELECT
    a.id AS actor_id,
    a.name,
    a.surname,
    a.patronymic,
    a.birthday,
    a.sex,
    a.information AS actor_information,
    f.id AS film_id,
    f.title,
    f.year,
    f.information AS film_information,
    f.rating
FROM
    actors a
        JOIN
    films_actors fa ON a.id = fa.actor_id
        JOIN
    films f ON fa.film_id = f.id;`

func (s *actorStorage) GetActorsWithFilms() ([]domain.ActorFilm, error) {
	rows, err := s.db.Queryx(getActorsWithFilms)
	if err != nil {
		return nil, err
	}
	actorsFilms := make(map[int64]domain.ActorFilm)

	for rows.Next() {
		var (
			actorId          int64
			name             string
			surname          string
			patronymic       sql.NullString
			birthday         time.Time
			sex              string
			actorInformation sql.NullString
			filmId           int64
			title            string
			year             int
			filmInformation  sql.NullString
			rating           sql.NullFloat64
		)

		err := rows.Scan(&actorId, &name, &surname, &patronymic, &birthday, &sex,
			&actorInformation, &filmId, &title, &year, &filmInformation, &rating)
		if err != nil {
			return nil, err
		}

		actor := domain.Actor{
			ID:          actorId,
			Name:        name,
			Surname:     surname,
			Patronymic:  patronymic,
			Birthday:    birthday,
			Sex:         sex,
			Information: actorInformation,
		}
		film := domain.Film{
			ID:          filmId,
			Title:       title,
			Year:        year,
			Information: filmInformation,
			Rating:      rating,
		}

		actorFilm, ok := actorsFilms[actor.ID]
		if !ok {
			actorFilm = domain.ActorFilm{Actor: actor}
		}
		actorFilm.Films = append(actorFilm.Films, film)
		actorsFilms[actor.ID] = actorFilm
	}
	var actorsFilmsArray []domain.ActorFilm

	for _, v := range actorsFilms {
		actorsFilmsArray = append(actorsFilmsArray, v)
	}

	return actorsFilmsArray, nil
}
