package main

import (
	"flag"
	"log"

	"github.com/4aykovski/tg-notion-bot/cmd/internal/client/telegram"
	zapLogger "github.com/4aykovski/tg-notion-bot/pkg/zap-logger"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	logger := zapLogger.New("development")

	token := mustToken()

	tgClient := telegram.New(tgBotHost, token)

	// fetcerh = fetcher.New(tgClient)

	// processor = fetcher.New(tgClient)

	// cosumer.Start(fetcher, processor)

}

func mustToken() string {
	token := flag.String(
		"bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
