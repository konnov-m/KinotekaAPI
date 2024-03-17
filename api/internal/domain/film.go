package domain

import "database/sql"

type Film struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Year        int             `json:"year"`
	Information sql.NullString  `json:"information"`
	Rating      sql.NullFloat64 `json:"rating"`
}

func (f *Film) IsValid() bool {
	return f.ID >= 0 && f.Year > 1000 && f.Title != ""
}
