package eventProcessor

import (
	"errors"
	"fmt"

	"github.com/4aykovski/tg-notion-bot/internal/client/notion"
	tgClient "github.com/4aykovski/tg-notion-bot/internal/client/telegram"
	"github.com/4aykovski/tg-notion-bot/internal/events"
	Logger "github.com/4aykovski/tg-notion-bot/pkg/logger"
)

var (
	ErrNoUpdates        = errors.New("fetch no updates")
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
	ErrUnknownDataType  = errors.New("unknown data type")
)

type speechAnalyzer interface {
	SpeechRecognizeOgg(fileName string) (text string, err error)
}

type aiBot interface {
	Completions(text string) (result string, err error)
}

type Processor struct {
	tg             *tgClient.Client
	speechAnalyzer speechAnalyzer
	aiBot          aiBot
	not            *notion.Client
	offset         int
	logger         *Logger.Logger
}

type Meta struct {
	ChatID   int
	Username string
}

type Data struct {
	Text  string
	Voice tgClient.Voice
}

func New(
	aiBot aiBot,
	speechAnalyzer speechAnalyzer,
	telegramClient *tgClient.Client,
	notionClient *notion.Client,
) *Processor {
	return &Processor{
		tg:             telegramClient,
		speechAnalyzer: speechAnalyzer,
		aiBot:          aiBot,
		not:            notionClient,
		logger:         Logger.New(),
	}
}
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("can't get events: %w", ErrNoUpdates)
	}

	res := make([]events.Event, 0, len(updates))

	for _, update := range updates {
		res = append(res, event(update))
	}

	p.offset = updates[len(updates)-1].Id + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("can't process event: %w", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process event: %w", err)
	}
	data, err := data(event)
	if err != nil {
		return fmt.Errorf("can't process event: %w", err)
	}

	if data.Voice != (tgClient.Voice{}) {
		if err := p.doCmdIfVoice(data.Voice, meta.ChatID); err != nil {
			return fmt.Errorf("can't process event: %w", err)
		}
	} else {
		if err := p.doCmdIfText(data.Text, meta.ChatID, meta.Username); err != nil {
			return fmt.Errorf("can't process event: %w", err)
		}
	}

	return nil
}

func data(event events.Event) (Data, error) {
	res, ok := event.Data.(Data)
	if !ok {
		return Data{}, fmt.Errorf("can't get data: %w", ErrUnknownDataType)
	}

	return res, nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta: %w", ErrUnknownMetaType)
	}

	return res, nil
}

func event(update tgClient.Update) events.Event {
	updateType := fetchType(update)

	res := events.Event{
		Type: updateType,
		Data: fetchData(update),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchData(update tgClient.Update) Data {
	if update.Message == nil {
		return Data{}
	}

	return Data{
		Text:  update.Message.Text,
		Voice: update.Message.Voice,
	}
}

func fetchType(update tgClient.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}

	return events.Message
}
