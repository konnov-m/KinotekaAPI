package storage

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"kinoteka/internal/domain"
)

type userStorage struct {
	db *sqlx.DB
}

func NewUserStorage(conn *sqlx.DB) UserStorage {
	return &userStorage{
		db: conn,
	}
}

const createUser = `INSERT INTO users (login, password) VALUES ($1, $2)`
const getRole = `SELECT id, name FROM roles WHERE name = $1`
const createUserRole = `INSERT INTO users_roles(user_id, role_id) VALUES ($1, $2)`

func (s *userStorage) CreateUser(user domain.User, role string) error {
	var r domain.Role
	err := s.db.Get(&r, getRole, role)
	if err != nil {
		return errors.New("There is no role")
	}

	_, err = s.db.Exec(createUser, user.Login, user.Password)
	if err != nil {
		return err
	}
	createdUser, err := s.GetUser(user.Login, user.Password)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(createUserRole, createdUser.ID, r.ID)

	return err
}

const getUserByLoginAndPass = `SELECT id, login, password FROM users WHERE login = $1 and password = $2`

func (s *userStorage) GetUser(login, pass string) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, getUserByLoginAndPass, login, pass)

	return &user, err
}

const getRoleByUserId = `SELECT role_id FROM users_roles WHERE user_id = $1`
const getRoleById = `SELECT id, name FROM roles WHERE id = $1`

func (s *userStorage) GetRole(userId int64) ([]domain.Role, error) {
	var rolesId []int64
	err := s.db.Select(&rolesId, getRoleByUserId, userId)
	if err != nil {
		return nil, err
	}
	roles := make([]domain.Role, 0)
	for _, roleId := range rolesId {
		var role domain.Role
		err := s.db.Get(&role, getRoleById, roleId)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
