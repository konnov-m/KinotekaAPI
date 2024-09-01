package service

import (
	"errors"
	"kinoteka/internal/domain"
	"kinoteka/internal/storage"
)

type filmService struct {
	s storage.FilmStorage
}

func NewFilmService(s storage.FilmStorage) Film {
	return &filmService{
		s: s,
	}
}

func (f *filmService) GetFilmsSortLike(orderBy, title string, desc bool) ([]domain.Film, error) {
	if orderBy != "title" && orderBy != "year" {
		orderBy = "rating"
	}

	return f.s.GetFilmsSortLike(orderBy, title, desc)
}

func (f *filmService) GetFilmsLike(title string) ([]domain.Film, error) {
	return f.s.GetFilmsLike(title)
}

func (f *filmService) GetFilmsSort(orderBy string, desc bool) ([]domain.Film, error) {
	if orderBy != "title" && orderBy != "year" {
		orderBy = "rating"
	}

	return f.s.GetFilmsSort(orderBy, desc)
}

func (f *filmService) GetFilm(id int64) (domain.Film, error) {
	return f.s.GetFilm(id)
}

func (f *filmService) CreateFilm(a domain.Film) error {
	if !a.IsValid() {
		return errors.New("film is not valid")
	}
	return f.s.CreateFilm(a)
}

func (f *filmService) UpdateFilm(a domain.Film) error {
	if !a.IsValid() {
		return errors.New("film is not valid")
	}
	return f.s.UpdateFilm(a)
}

func (f *filmService) DeleteFilm(id int64) error {
	if err := f.s.DeleteFilmsActors(id); err != nil {
		return err
	}

	return f.s.DeleteFilm(id)
}

func (f *filmService) SearchFilmsWithActor(substr string) ([]domain.ActorFilm, error) {
	return f.s.SearchFilmsWithActor(substr)
}

func (f *filmService) AddActorToFilm(filmId int64, actorId []int64) error {
	return f.s.AddActorToFilm(filmId, actorId)
}
