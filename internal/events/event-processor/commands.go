package eventProcessor

import (
	"path/filepath"
	"strings"

	"github.com/4aykovski/tg-notion-bot/config"
	"github.com/4aykovski/tg-notion-bot/internal/client/telegram"
	"github.com/4aykovski/tg-notion-bot/lib/helpers"
	"go.uber.org/zap"
)

func (p *Processor) doCmdIfVoice(voice telegram.Voice, chatId int, username string) (err error) {
	defer func() { err = helpers.ErrWrapIfNotNil("can't do command with voice", err) }()
	f, err := p.tg.FileInfo(voice.FileId)
	if err != nil {
		return err
	}

	if err := p.tg.DownloadFile(f.FilePath); err != nil {
		return err
	}

	voiceText, err := p.speechAnalyzer.SpeechRecognizeOgg(fileName(f.FilePath))
	if err != nil {
		return err
	}

	editedText, err := p.aiBot.Completions(voiceText)
	if err != nil {
		return err
	}

	err = p.not.CreateNewPageInDatabase(config.NotionDatabaseId, editedText)
	if err != nil {
		return err
	}

	err = p.tg.SendMessage(chatId, msgSuccessfulSaved)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) doCmdIfText(text string, chatID int, username string) (err error) {
	defer func() { err = helpers.ErrWrapIfNotNil("can't do command with text", err) }()

	text = strings.TrimSpace(text)

	p.logger.Info("got new command",
		zap.String("username", username),
		zap.Int("chatID", chatID))

	switch text {
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func fileName(filePath string) string {
	_, fName := filepath.Split(filePath)
	return fName
}
