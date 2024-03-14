package storage

import (
	"KinotekaAPI/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type FilmStorage interface {
	GetFilmsLike(title string) ([]domain.Film, error)
	GetFilmsSort(orderBy string, desc bool) ([]domain.Film, error)
	GetFilmsSortLike(orderBy, title string, desc bool) ([]domain.Film, error)
	GetFilm(id int64) (domain.Film, error)
	CreateFilm(a domain.Film) error
	UpdateFilm(a domain.Film) error
}

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
	if err != nil {
		return nil, err
	}

	return films, nil
}

const getFilmsSort = `SELECT id, title, year, information, rating FROM films ORDER BY`

func (s *filmStorage) GetFilmsSort(orderBy string, desc bool) ([]domain.Film, error) {
	if orderBy != "title" && orderBy != "year" {
		orderBy = "rating"
	}

	var sql string

	if desc {
		sql = fmt.Sprintf("%s %s %s", getFilmsSort, orderBy, "DESC")
	} else {
		sql = fmt.Sprintf("%s %s %s", getFilmsSort, orderBy, "ASC")
	}

	var films []domain.Film
	err := s.db.Select(&films, sql)
	if err != nil {
		return nil, err
	}

	return films, nil
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
	if err != nil {
		return nil, err
	}

	return films, nil
}

const getFilmId = `SELECT id, title, year, information, rating
FROM films WHERE id = $1`

func (s *filmStorage) GetFilm(id int64) (domain.Film, error) {
	var film domain.Film
	err := s.db.Get(&film, getFilmId, id)
	if err != nil {
		return domain.Film{}, err
	}

	return film, nil
}

const saveFilm = `INSERT INTO films (title, year, information, rating)
VALUES ($1, $2, $3, $4);`

func (s *filmStorage) CreateFilm(a domain.Film) error {
	_, err := s.db.Exec(saveFilm, a.Title, a.Year, a.Information, a.Rating)
	if err != nil {
		return err
	}

	return nil
}

const updateFilm = `UPDATE films SET title=$1, year=$2, information=$3, rating=$4 WHERE id=$5;`

func (s *filmStorage) UpdateFilm(a domain.Film) error {
	_, err := s.db.Exec(updateFilm, a.Title, a.Year, a.Information, a.Rating, a.ID)
	if err != nil {
		return err
	}

	return nil
}
