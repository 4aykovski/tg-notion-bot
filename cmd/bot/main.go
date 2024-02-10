package main

import (
	"github.com/4aykovski/tg-notion-bot/config"
	gigachatClient "github.com/4aykovski/tg-notion-bot/internal/client/gigachat"
	"github.com/4aykovski/tg-notion-bot/internal/client/notion"
	salutespeechClient "github.com/4aykovski/tg-notion-bot/internal/client/salutespeech"
	telegramClient "github.com/4aykovski/tg-notion-bot/internal/client/telegram"
	"github.com/4aykovski/tg-notion-bot/internal/consumer/eventConsumer"
	eventProcessor "github.com/4aykovski/tg-notion-bot/internal/events/event-processor"
	Logger "github.com/4aykovski/tg-notion-bot/pkg/logger"
)

func main() {
	logger := Logger.New(config.Type)

	tgClient, err := telegramClient.New(config.TgBotHost, config.TgBotToken)
	if err != nil {
		logger.Fatal(err.Error())
	}

	spClient, err := salutespeechClient.New(config.SalutespeechToken)
	if err != nil {
		logger.Fatal(err.Error())
	}

	gcClient, err := gigachatClient.New(config.GigaChatToken)
	if err != nil {
		logger.Fatal(err.Error())
	}

	notClient, err := notion.New(config.NotionIntegrationToken)
	if err != nil {
		logger.Fatal(err.Error())
	}

	eP := eventProcessor.New(gcClient, spClient, tgClient, notClient)

	logger.Info("service started")

	consumer := eventConsumer.New(eP, eP, config.BatchSize)

	if err := consumer.Start(); err != nil {
		logger.Fatal(err.Error())
	}

}
