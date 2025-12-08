package dbstore

import (
	"TODO/adapter/driven/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Auth struct {
	db *sqlx.DB
}

// Login проверяет пользователя в базе
func (r *Auth) Login(username, password string) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, "SELECT id, username, password FROM users WHERE username=$1 AND password=$2", username, password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return &u, nil
}

// Stats возвращает количество пользователей
func (r *Auth) Stats() (map[string]interface{}, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users")
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"users_count": count,
	}, nil
}
