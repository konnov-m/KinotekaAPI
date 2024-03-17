package domain

type User struct {
	ID       int64
	Login    string
	Password string
} // @name User

type Role struct {
	ID   int64
	Name string
} // @name Role
