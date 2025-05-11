package subpub

import (
	"SubPub_project/pkg/logger"

	"context"
	"errors"
	"sync"
)

type Publisher struct {
	mu       sync.RWMutex
	subjects map[string][]*Subscriber
	closed   bool
	closeCh  chan struct{}
	log      logger.Logger
}

type PublisherInterface interface {
	Publish(ctx context.Context, subject string, msg interface{}) error
}

func NewPublisher(ctx context.Context, log logger.Logger) *Publisher {
	p := &Publisher{
		subjects: make(map[string][]*Subscriber),
		closeCh:  make(chan struct{}),
		log:      log,
	}
	go func() {
		<-ctx.Done()
		p.log.Info("Закрываем `Publisher` по контексту")
		p.mu.Lock()
		p.closed = true
		close(p.closeCh)
		p.mu.Unlock()
	}()

	return p
}

func (p *Publisher) Publish(ctx context.Context, subject string, msg interface{}) error {

	p.log.Info("Публикация вызвана:", subject, "сообщение: ", msg)

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return errors.New("Система закрыта, публикация невозможна")
	}

	subs := append([]*Subscriber{}, p.subjects[subject]...)
	if len(subs) == 0 {
		p.log.Info("Нет подписчиков для: %s", subject)
	}

	for _, sub := range subs {
		p.log.Info("Передача сообщения подписчику:", msg)

		select {
		case sub.ch <- msg:
			p.log.Info("Сообщение отправлено подписчику:", sub)
		case <-ctx.Done():
			return ctx.Err()
		default:
			p.log.Info("Подписчик не смог обработать сообщение, для подписчика", sub)
		}
	}
	return nil
}
