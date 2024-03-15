package storage

import (
	"KinotekaAPI/internal/domain"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type filmStorage struct {
	db *sqlx.DB
}

func NewFilmStorage(conn *sqlx.DB) FilmStorage {
	return &filmStorage{
		db: conn,
	}
}

const getFilmsLike = `SELECT * FROM films WHERE LOWER(title) LIKE '%' || LOWER($1) || '%';`

func (s *filmStorage) GetFilmsLike(title string) ([]domain.Film, error) {

	var films []domain.Film
	err := s.db.Select(&films, getFilmsLike, title)

	return films, err
}

const getFilmsSort = `SELECT id, title, year, information, rating FROM films ORDER BY`

func (s *filmStorage) GetFilmsSort(orderBy string, desc bool) ([]domain.Film, error) {
	var sql string

	if desc {
		sql = fmt.Sprintf("%s %s %s", getFilmsSort, orderBy, "DESC")
	} else {
		sql = fmt.Sprintf("%s %s %s", getFilmsSort, orderBy, "ASC")
	}

	var films []domain.Film
	err := s.db.Select(&films, sql)

	return films, err
}

const getFilmsSortLike = `SELECT * FROM films WHERE LOWER(title) LIKE '%' || LOWER($1) || '%' ORDER BY`

func (s *filmStorage) GetFilmsSortLike(orderBy, title string, desc bool) ([]domain.Film, error) {
	if orderBy != "title" && orderBy != "year" {
		orderBy = "rating"
	}

	var sql string

	if desc {
		sql = fmt.Sprintf("%s %s %s", getFilmsSortLike, orderBy, "DESC")
	} else {
		sql = fmt.Sprintf("%s %s %s", getFilmsSortLike, orderBy, "ASC")
	}

	var films []domain.Film
	err := s.db.Select(&films, sql, title)

	return films, err
}

const getFilmId = `SELECT id, title, year, information, rating
FROM films WHERE id = $1`

func (s *filmStorage) GetFilm(id int64) (domain.Film, error) {
	var film domain.Film
	err := s.db.Get(&film, getFilmId, id)

	return film, err
}

const saveFilm = `INSERT INTO films (title, year, information, rating)
VALUES ($1, $2, $3, $4);`

func (s *filmStorage) CreateFilm(a domain.Film) error {
	_, err := s.db.Exec(saveFilm, a.Title, a.Year, a.Information, a.Rating)

	return err
}

const updateFilm = `UPDATE films SET title=$1, year=$2, information=$3, rating=$4 WHERE id=$5;`

func (s *filmStorage) UpdateFilm(a domain.Film) error {
	_, err := s.db.Exec(updateFilm, a.Title, a.Year, a.Information, a.Rating, a.ID)

	return err
}

const deleteFilm = `DELETE from films WHERE id=$1;`

func (s *filmStorage) DeleteFilm(id int64) error {
	_, err := s.db.Exec(deleteFilm, id)

	return err
}

const searchFilmsWithActor = `
SELECT
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
    films f ON fa.film_id = f.id
WHERE LOWER(name) like '%' || $1 || '%' or
      LOWER(surname) like '%' || $1 || '%' or
      LOWER(patronymic) like '%' || $1 || '%';
`

func (s *filmStorage) SearchFilmsWithActor(substr string) ([]domain.ActorFilm, error) {
	rows, err := s.db.Queryx(searchFilmsWithActor, strings.ToLower(substr))
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

const deleteFilmsActors = `DELETE FROM films_actors WHERE film_id = $1`

func (s *filmStorage) DeleteFilmsActors(id int64) error {
	_, err := s.db.Exec(deleteFilmsActors, id)

	return err
}
