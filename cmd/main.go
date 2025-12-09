package main

import (
	"TODO/adapter/driven/dbstore"

	"TODO/adapter/driving/telegram"
	"TODO/config"
	"TODO/pkg/logger"
	"TODO/services"
	"context"
	"time"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Подключение к БД
	pg, err := sqlx.Connect("postgres", cfg.Postgres.ConnectionURL())
	if err != nil {
		log.Error().Err(err).Msg("Ошибка подключения к БД")
		return
	}
	defer pg.Close()

	// Telegram
	bot, err := tgbot.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка подключения к Telegram")
		return
	}

	// Репозиторий и сервис
	todoRepo := dbstore.NewRepo(pg)
	todoSvc := services.NewService(todoRepo)

	// Telegram Router
	tgRouter := telegram.NewRouter(bot, todoSvc)
	tgRouter.StartDueTaskNotifier(ctx, 1*time.Hour)

	// Запуск Telegram бота
	go func() {
		log.Info().Msg("Telegram bot started")
		if err := tgRouter.Run(ctx); err != nil {
			log.Error().Err(err).Msg("Telegram bot error")
		}
	}()

	// Ожидание завершения
	<-ctx.Done()
	log.Info().Msg("Shutting down")
}
