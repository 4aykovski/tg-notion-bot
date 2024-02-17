package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database            DatabaseConfig
	Telegram            TelegramConfig
	Notion              NotionConfig
	GigaChat            GigaChatConfig
	Salutespeech        SalutespeechConfig
	VoicesFileDirectory string
	BatchSize           int
}

type DatabaseConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	DbName      string
	DSNTemplate string
}

type TelegramConfig struct {
	Host  string
	Token string
}

type NotionConfig struct {
	Host             string
	APIBasePath      string
	Version          string
	IntegrationToken string
}

type GigaChatConfig struct {
	Host        string
	APIBasePath string
	Token       string
	Auth        string
}

type SalutespeechConfig struct {
	Host        string
	APIBasePath string
	Token       string
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("can't load .env file: %w", err)
	}

	postgresPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return nil, fmt.Errorf("can't get postgresPort: %w", err)
	}

	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		batchSize = 100 // default
	}

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     postgresPort,
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRESQL_CHAYKOVSKI_PASSWORD"),
			DbName:   os.Getenv("POSTGRES_DB_NAME"),
		},
		Telegram: TelegramConfig{
			Host:  "api.telegram.org",
			Token: os.Getenv("TG_BOT_TOKEN"),
		},
		Notion: NotionConfig{
			Host:             "api.notion.com",
			APIBasePath:      "/v1",
			Version:          "2022-06-28",
			IntegrationToken: os.Getenv("NOTION_INTEGRATION_TOKEN"),
		},
		GigaChat: GigaChatConfig{
			Host:        "gigachat.devices.sberbank.ru",
			APIBasePath: "/api/v1/",
			Token:       os.Getenv("GIGACHAT_ACCESS_TOKEN"),
			Auth:        os.Getenv("GIGACHAT_AUTH_TOKEN"),
		},
		Salutespeech: SalutespeechConfig{
			Host:        "smartspeech.sber.ru",
			APIBasePath: "/rest/v1/",
			Token:       os.Getenv("SALUTESPEECH_ACCESS_TOKEN"),
		},
		VoicesFileDirectory: "./voices/",
		BatchSize:           batchSize,
	}

	cfg.Database.DSNTemplate = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName)

	return cfg, nil
}
