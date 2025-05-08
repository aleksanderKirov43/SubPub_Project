package subpub

import (
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

func NewPublisher() *Publisher {
	return &Publisher{
		subjects: make(map[string][]*Subscriber),
		closeCh:  make(chan struct{}),
	}
}

func (p *Publisher) Publish(subject string, msg interface{}) error {
	log.Println("Публикация вызвана:", subject, msg)
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.closed {
		return errors.New("Система закрыта")
	}

	subs := append([]*Subscriber{}, p.subjects[subject]...)
	if len(subs) == 0 {
		log.Println("Нет подписчиков для:", subject)
	}

	for _, sub := range subs {
		log.Println("Передача в канал подписчика:", msg)
		select {
		case sub.ch <- msg:
			log.Println("Сообщение отправлено подписчику:", sub)
		default:
			log.Println("Подписчик не смог обработать сообщение!")
		}
	}
	return nil
}
