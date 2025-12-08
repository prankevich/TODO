package usecase

import (
	"TODO/adapter/driven/dbstore"
	"TODO/config"
	authenticate "TODO/usecase/authenticator"
	"TODO/usecase/emails_getter"
	usercreater "TODO/usecase/user_creater"
)

authenticate

	usercreater
)

type UseCases struct {
	UserCreater   usecase.UserCreater
	Authenticator usecase.Authenticate
	EmailsGetter  usecase.EmailsGetter
}

func New(cfg config.Config,
	store *dbstore.DBStore,
	amqp driven.AmqpProducer) *UseCases {
	return &UseCases{
		UserCreater:   usercreater.New(&cfg, store.UserStorage, amqp),
		Authenticator: authenticate.New(&cfg, store.UserStorage),
		EmailsGetter:  emails_getter.New(&cfg, store.UserStorage),
	}
}
