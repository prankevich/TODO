package usecase

import (
	"TODO/adapter/driven/models"
	"context"
)

type Authenticate interface {
	Authenticate(ctx context.Context, user models.User) (int, error)
}
