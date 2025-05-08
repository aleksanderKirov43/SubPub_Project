package subpub

import (
	"context"
	"errors"
	"log"
	"sync"
)

type Publisher struct {
	mu       sync.RWMutex
	subjects map[string][]*Subscriber
	closed   bool
	closeCh  chan struct{}
}

type PublisherInterface interface {
	Publish(ctx context.Context, subject string, msg interface{}) error
}

func NewPublisher(ctx context.Context) *Publisher {
	p := &Publisher{
		subjects: make(map[string][]*Subscriber),
		closeCh:  make(chan struct{}),
	}
	go func() {
		select {
		case <-ctx.Done():
			log.Println("Закрываем `Publisher` по контексту")
			p.mu.Lock()
			p.closed = true
			close(p.closeCh)
			p.mu.Unlock()
		}
	}()

	return p
}

func (p *Publisher) Publish(ctx context.Context, subject string, msg interface{}) error {

	log.Println("Публикация вызвана:", subject, "сообщение: ", msg)

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return errors.New("Система закрыта, публикация невозможна")
	}

	subs := append([]*Subscriber{}, p.subjects[subject]...)
	if len(subs) == 0 {
		log.Println("Нет подписчиков для:", subject)
	}

	for _, sub := range subs {
		log.Println("Передача сообщения подписчику:", msg)

		select {
		case sub.ch <- msg:
			log.Println("Сообщение отправлено подписчику:", sub)
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Println("Подписчик не смог обработать сообщение!")
		}
	}
	return nil
}
