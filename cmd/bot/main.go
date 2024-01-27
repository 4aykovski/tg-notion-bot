package main

import (
	"flag"
	"log"

	zapLogger "github.com/4aykovski/tg-notion-bot/pkg/zap-logger"
)

func main() {
	logger := zapLogger.New("development")

	token := mustToken()

	logger.Info(token)

	// tgClient = telegram.New(token)

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
