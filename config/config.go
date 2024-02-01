package config

import (
	"fmt"
	"os"
)

const (
	Type = "development"
)

const (
	TgBotHost                = "api.telegram.org"
	NotionHost               = "api.notion.com"
	NotionAPIBasePath        = "/v1"
	NotionVersion            = "2022-06-28"
	NotionDatabaseId         = "42bc30c24c184db0a890018c009b69fd"
	SalutespeechHost         = "smartspeech.sber.ru"
	SalutesleepchAPIBasePath = "/rest/v1/"
	GigaChatHost             = "gigachat.devices.sberbank.ru"
	GigaChatAPIBasePath      = "/api/v1/"
	BatchSize                = 100
	VoicesFileDirectory      = "./voices/"
)

var (
	TgBotToken             = os.Getenv("TG_BOT_TOKEN")
	NotionIntegrationToken = os.Getenv("NOTION_INTEGRATION_TOKEN")
	SalutespeechToken      = os.Getenv("SALUTESPEECH_ACCESS_TOKEN")
	GigaChatToken          = os.Getenv("GIGACHAT_ACCESS_TOKEN")
	PostgresHost           = "localhost"
	PostgresPort           = 5432
	PostgresUser           = "chaykovski"
	PostgresPassword       = os.Getenv("POSTGRESQL_CHAYKOVSKI_PASSWORD")
	PostgresDbName         = "test"
	PostgresDSN            = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		PostgresHost, PostgresPort, PostgresUser, PostgresPassword, PostgresDbName)
)
