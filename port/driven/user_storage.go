package driven

import (
	"TODO/adapter/driven/models"
	"context"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user models.User) (err error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	GetAllUsersEmails(ctx context.Context) (emails []string, err error)
}
