package main

import (
	"github.com/4aykovski/tg-notion-bot/config"
	gigachatClient "github.com/4aykovski/tg-notion-bot/internal/client/gigachat"
	"github.com/4aykovski/tg-notion-bot/internal/client/notion"
	salutespeechClient "github.com/4aykovski/tg-notion-bot/internal/client/salutespeech"
	telegramClient "github.com/4aykovski/tg-notion-bot/internal/client/telegram"
	"github.com/4aykovski/tg-notion-bot/internal/consumer/event"
	eventProcessor "github.com/4aykovski/tg-notion-bot/internal/events/processor"
	"github.com/4aykovski/tg-notion-bot/internal/storage/postgres"
	Logger "github.com/4aykovski/tg-notion-bot/pkg/logger"
)

func main() {
	logger := Logger.New()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err.Error())
	}

	tgClient, err := telegramClient.New(cfg.Telegram, cfg.VoicesFileDirectory)
	if err != nil {
		logger.Fatal(err.Error())
	}

	spClient, err := salutespeechClient.New(cfg.Salutespeech, cfg.VoicesFileDirectory)
	if err != nil {
		logger.Fatal(err.Error())
	}

	gcClient, err := gigachatClient.New(cfg.GigaChat)
	if err != nil {
		logger.Fatal(err.Error())
	}

	notClient, err := notion.New(cfg.Notion)
	if err != nil {
		logger.Fatal(err.Error())
	}

	psqlDatabase, err := postgres.NewPostgresDatabase(cfg.Database)
	if err != nil {
		logger.Fatal(err.Error())
	}

	userRepository := postgres.NewUserRepository(psqlDatabase)

	eP := eventProcessor.New(gcClient, spClient, tgClient, notClient, userRepository)

	logger.Info("service started")

	consumer := event.New(eP, eP, cfg.BatchSize)

	if err := consumer.Start(); err != nil {
		logger.Fatal(err.Error())
	}
}
