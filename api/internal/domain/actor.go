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
}
