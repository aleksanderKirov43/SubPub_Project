package subpub

import (
	"context"
	"fmt"
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

var ErrPubSubClosed = fmt.Errorf("система подписок закрыта")

func (s *SubPub) Subscribe(ctx context.Context, subject string, cb MessageHandler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, ErrPubSubClosed
	}

	sub := NewSubscriber(ctx, cb)
	s.subscribers[subject] = append(s.subscribers[subject], sub)

	go func() {
		//<-ctx.Done()
		sub.Unsubscribe()
	}()

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
		return nil
	}

	for _, sub := range subs {
		select {
		case sub.ch <- msg:
			fmt.Println("Сообщение доставлено", sub)
		default:
			fmt.Println("Подписчик не успел обработать сообщение!")
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

	s.subscribers = make(map[string][]*Subscriber)
	fmt.Println("Система подписок закрыта")
	return nil
}
