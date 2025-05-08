package subpub

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type SubPubInterface interface {
	Subscribe(ctx context.Context, subject string, cb MessageHandler) (Subscription, error)
	Publish(ctx context.Context, subject string, msg interface{}) error
	Close(ctx context.Context) error
}

type SubPub struct {
	mu          sync.RWMutex
	subscribers map[string][]*Subscriber
	closed      bool
}

func NewSubPub() SubPubInterface {
	return &SubPub{
		subscribers: make(map[string][]*Subscriber),
		closed:      false,
	}
}

func (s *SubPub) Subscribe(ctx context.Context, subject string, cb MessageHandler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, ErrPubSubClosed
	}

	sub := NewSubscriber(ctx, cb)
	s.subscribers[subject] = append(s.subscribers[subject], sub)

	go func() {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
		}
	}()

	log.Println("Подписка зарегистрирована для ключа:", subject)

	return sub, nil
}

func (s *SubPub) Publish(ctx context.Context, subject string, msg interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return ErrPubSubClosed
	}

	subs := s.subscribers[subject]
	if len(subs) == 0 {
		log.Println("Нет подписчиков для:", subject)
		return nil
	}

	for _, sub := range subs {
		log.Println("Отправка сообщения подписчику:", msg)
		select {
		case sub.ch <- msg:
			log.Println("Сообщение доставлено")
		default:
			log.Println("Подписчик не успел обработать сообщение!")
		}
	}
	return nil
}

func (s *SubPub) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	for _, subs := range s.subscribers {
		for _, sub := range subs {
			sub.Unsubscribe()
		}
	}

	s.subscribers = make(map[string][]*Subscriber) // ✅ Очищаем подписчиков после закрытия
	log.Println("Система подписок закрыта")
	return nil
}

var ErrPubSubClosed = fmt.Errorf("система подписок закрыта")
