package eventConsumer

import (
	"fmt"
	"sync"
	"time"

	"github.com/4aykovski/tg-notion-bot/internal/events"
	Logger "github.com/4aykovski/tg-notion-bot/pkg/logger"
	"go.uber.org/zap"
)

type Fetcher interface {
	Fetch(limit int) ([]events.Event, error)
}

type Processor interface {
	Process(e events.Event) error
}

type Consumer struct {
	fetcher   Fetcher
	processor Processor
	batchSize int
	logger    *Logger.Logger
}

func New(fetcher Fetcher, processor Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		logger:    Logger.New(),
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			c.logger.Error("can't handle event",
				zap.String("error", err.Error()))

			continue
		}
	}
}

func (c *Consumer) handleEvents(fetchedEvents []events.Event) error {
	wg := sync.WaitGroup{}
	wg.Add(len(fetchedEvents))
	for _, event := range fetchedEvents {
		go func(event events.Event) {
			defer wg.Done()
			data := fmt.Sprintf("%v", event.Data)
			c.logger.Info("got new event",
				zap.String("data", data),
				zap.Int("type", int(event.Type)))

			if err := c.processor.Process(event); err != nil {
				data := fmt.Sprintf("%v", event.Data)
				c.logger.Error("can't handle event",
					zap.String("error", err.Error()),
					zap.String("data", data),
					zap.Int("type", int(event.Type)))
			}
		}(event)
	}
	wg.Wait()

	return nil
}
