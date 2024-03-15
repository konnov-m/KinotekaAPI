package domain

type User struct {
	ID       int64
	Login    string
	Password string
}

type Role struct {
	ID   int64
	Name string
}
