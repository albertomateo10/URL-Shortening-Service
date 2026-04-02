package service

import (
	"context"
	"log"
	"time"

	"github.com/albertomateo10/url-shortener/backend/internal/model"
	"github.com/albertomateo10/url-shortener/backend/internal/repository"
)

const (
	clickChannelSize = 1000
	flushInterval    = 1 * time.Second
	maxBatchSize     = 100
)

type ClickLogger struct {
	clickRepo *repository.ClickRepository
	urlRepo   *repository.URLRepository
	eventCh   chan *model.ClickEvent
}

func NewClickLogger(clickRepo *repository.ClickRepository, urlRepo *repository.URLRepository) *ClickLogger {
	cl := &ClickLogger{
		clickRepo: clickRepo,
		urlRepo:   urlRepo,
		eventCh:   make(chan *model.ClickEvent, clickChannelSize),
	}
	go cl.processEvents()
	return cl
}

func (cl *ClickLogger) Log(event *model.ClickEvent) {
	select {
	case cl.eventCh <- event:
	default:
		log.Println("click logger: channel full, dropping event")
	}
}

func (cl *ClickLogger) processEvents() {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	batch := make([]*model.ClickEvent, 0, maxBatchSize)

	for {
		select {
		case event := <-cl.eventCh:
			batch = append(batch, event)
			if len(batch) >= maxBatchSize {
				cl.flush(batch)
				batch = make([]*model.ClickEvent, 0, maxBatchSize)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				cl.flush(batch)
				batch = make([]*model.ClickEvent, 0, maxBatchSize)
			}
		}
	}
}

func (cl *ClickLogger) flush(batch []*model.ClickEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, event := range batch {
		if err := cl.clickRepo.Insert(ctx, event); err != nil {
			log.Printf("click logger: insert error: %v", err)
			continue
		}
		if err := cl.urlRepo.IncrementClickCount(ctx, event.URLID); err != nil {
			log.Printf("click logger: increment click count error: %v", err)
		}
	}
}
