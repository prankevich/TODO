package usercreater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type UseCase struct {
	cfg         *config.Config
	userStorage driven.UserStorage
	amqp        driven.AmqpProducer
}

func New(cfg *config.Config, userStorage driven.UserStorage, amqp driven.AmqpProducer) *UseCase {
	return &UseCase{
		cfg:         cfg,
		userStorage: userStorage,
		amqp:        amqp,
	}
}

func (u *UseCase) CreateUser(ctx context.Context, user domain.User) (err error) {
	// Проверить существует ли пользователь с таким username'ом в бд
	_, err = u.userStorage.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotfound) {
			return err
		}
	} else {
		return errs.ErrUsernameAlreadyExists
	}

	// За хэшировать пароль
	user.Password, err = utils.GenerateHash(user.Password)
	if err != nil {
		return err
	}

	user.Role = domain.RoleUser

	// Добавить пользователя в бд
	if err = u.userStorage.CreateUser(ctx, user); err != nil {
		return err
	}

	go func() {
		msg := amqp.Message{
			Recipient: user.Email,
		}

		rawMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("Error marshalling message: %v\n", err)
		}

		if err = u.amqp.Publish(ctx, "auth-queue", rawMsg); err != nil {
			fmt.Printf("Error publishing message: %v\n", err)
		}
	}()

	return nil
}
