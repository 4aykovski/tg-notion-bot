package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/4aykovski/tg-notion-bot/internal/client/telegram"
	"go.uber.org/zap"
)

func (p *Processor) doCmdIfVoice(voice telegram.Voice, chatId int) (err error) {
	f, err := p.tg.FileInfo(voice.FileId)
	if err != nil {
		return fmt.Errorf("can't do command with voice: %w", err)
	}

	if err := p.tg.DownloadFile(f.FilePath); err != nil {
		return fmt.Errorf("can't do command with voice: %w", err)
	}

	voiceText, err := p.speechAnalyzer.SpeechRecognizeOgg(fileName(f.FilePath))
	if err != nil {
		return fmt.Errorf("can't do command with voice: %w", err)
	}

	editedText, err := p.aiBot.Completions(voiceText)
	if err != nil {
		return fmt.Errorf("can't do command with voice: %w", err)
	}

	notionDatabaseId := "42bc30c24c184db0a890018c009b69fd"

	err = p.notion.CreateNewPageInDatabase(notionDatabaseId, editedText)
	if err != nil {
		return fmt.Errorf("can't do command with voice: %w", err)
	}

	err = p.tg.SendMessage(chatId, msgSuccessfulSaved)
	if err != nil {
		return fmt.Errorf("can't do command with voice: %w", err)
	}

	return nil
}

func (p *Processor) doCmdIfText(text string, chatID int, username string) (err error) {
	text = strings.TrimSpace(text)

	p.logger.Info("got new command",
		zap.String("username", username),
		zap.Int("chatID", chatID))

	switch text {
	default:
		err := p.tg.SendMessage(chatID, msgUnknownCommand)
		if err != nil {
			return fmt.Errorf("can't do command with text: %w", err)
		}
		return nil
	}
}

func fileName(filePath string) string {
	_, fName := filepath.Split(filePath)
	return fName
}
