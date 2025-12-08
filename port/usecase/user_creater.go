package usecase

import (
	"TODO/adapter/driven/models"
	"context"
)

type UserCreater interface {
	CreateUser(ctx context.Context, user models.User) (err error)
}
