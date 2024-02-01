package eventProcessor

import (
	"errors"

	"github.com/4aykovski/tg-notion-bot/config"
	gcClient "github.com/4aykovski/tg-notion-bot/internal/client/gigachat"
	"github.com/4aykovski/tg-notion-bot/internal/client/notion"
	spClient "github.com/4aykovski/tg-notion-bot/internal/client/salutespeech"
	tgClient "github.com/4aykovski/tg-notion-bot/internal/client/telegram"
	"github.com/4aykovski/tg-notion-bot/internal/events"
	"github.com/4aykovski/tg-notion-bot/lib/helpers"
	zapLogger "github.com/4aykovski/tg-notion-bot/pkg/zap-logger"
)

var (
	ErrNoUpdates        = errors.New("fetch no updates")
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
	ErrUnknownDataType  = errors.New("unknown data type")
)

type Processor struct {
	tg     *tgClient.Client
	sp     *spClient.Client
	gc     *gcClient.Client
	not    *notion.Client
	offset int
	logger *zapLogger.Logger
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
	gigaChatClient *gcClient.Client,
	salutespeechClient *spClient.Client,
	telegramClient *tgClient.Client,
	notionClient *notion.Client,
) *Processor {
	return &Processor{
		tg:     telegramClient,
		sp:     salutespeechClient,
		gc:     gigaChatClient,
		not:    notionClient,
		logger: zapLogger.New(config.Type),
	}
}
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, helpers.ErrWrapIfNotNil("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, ErrNoUpdates
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
		return helpers.ErrWrapIfNotNil("can't process event", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return helpers.ErrWrapIfNotNil("can't process event", err)
	}
	data, err := data(event)
	if err != nil {
		return helpers.ErrWrapIfNotNil("can't process event", err)
	}

	if data.Voice != (tgClient.Voice{}) {
		if err := p.doCmdIfVoice(data.Voice, meta.ChatID, meta.Username); err != nil {
			return helpers.ErrWrapIfNotNil("can't process event", err)
		}
	} else {
		if err := p.doCmdIfText(data.Text, meta.ChatID, meta.Username); err != nil {
			return helpers.ErrWrapIfNotNil("can't process event", err)
		}
	}

	return nil
}

func data(event events.Event) (Data, error) {
	res, ok := event.Data.(Data)
	if !ok {
		return Data{}, helpers.ErrWrapIfNotNil("can't get meta", ErrUnknownDataType)
	}

	return res, nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, helpers.ErrWrapIfNotNil("can't get meta", ErrUnknownMetaType)
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
