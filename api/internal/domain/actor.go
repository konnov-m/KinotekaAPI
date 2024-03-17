package domain

import (
	"database/sql"
	"time"
)

type Actor struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Surname     string         `json:"surname"`
	Patronymic  sql.NullString `json:"patronymic"`
	Birthday    time.Time      `json:"birthday"`
	Sex         string         `json:"sex"`
	Information sql.NullString `json:"information"`
} // @name Actor

type ActorFilm struct {
	Actor Actor  `db:"actors"`
	Films []Film `db:"films"`
} // @name ActorFilm

func (a *Actor) IsValid() bool {
	return a.ID >= 0 && a.Name != "" && a.Surname != "" && a.Sex != "" && !a.Birthday.IsZero()
}
