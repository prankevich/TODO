package emails_getter

import (
	"TODO/config"
	"context"
	"fmt"
)

type UseCase struct {
	cfg         *config.Config
	userStorage driven.UserStorage
}

func New(cfg *config.Config, userStorage driven.UserStorage) *UseCase {
	return &UseCase{
		cfg:         cfg,
		userStorage: userStorage,
	}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]string, error) {
	emails, err := uc.userStorage.GetAllUsersEmails(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase.GetAllUsersEmails: %w", err)
	}

	return emails, nil
}
