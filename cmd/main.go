package main

import (
	"TODO/adapter/driven/dbstore"
	"TODO/adapter/driving/telegram"
	"TODO/config"
	"TODO/pkg/logger"
	"TODO/services"
	"context"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	_ = godotenv.Load(".env")
	var cfg config.Config
	if err := envconfig.ProcessWith(context.TODO(), &envconfig.Config{
		Target:   &cfg,
		Lookuper: envconfig.OsLookuper(),
	}); err != nil {
		panic(err)
	}
	log := logger.New()
	ctx := context.Background()

	//  Postgres
	pg, err := sqlx.Connect("postgres", cfg.Postgres.ConnectionURL())
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pg.Close()

	// Telegram
	bot, err := tgbot.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatalf("telegram: %v", err)
	}

	// Repos
	todoRepo := dbstore.NewRepo(pg)
	// Services
	todoSvc := services.NewService(todoRepo)
	// Router
	tgRouter := telegram.NewRouter(bot, todoSvc)

	go func() {
		log.Printf("telegram bot started")
		if err := tgRouter.Run(ctx); err != nil {
			log.Printf("telegram error: %v", err)
		}
	}()
	<-ctx.Done()
	log.Printf("bye")
}
