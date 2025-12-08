package usecase

import (
	"context"
)

type EmailsGetter interface {
	GetAll(ctx context.Context) ([]string, error)
}
