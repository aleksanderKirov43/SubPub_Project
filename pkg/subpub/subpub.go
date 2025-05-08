package subpub

import (
	"context"
	"sync"
)

type SubPubInterface interface {
	Subscribe(subject string, cb MessageHandler) (Subscription, error)
	Publish(subject string, msg interface{}) error
	Close(ctx context.Context) error
}

type SubPub struct {
	mu          sync.RWMutex
	publisher   *Publisher
	subscribers map[string][]*Subscriber
}

func NewSubPub() SubPubInterface {
	return &SubPub{
		publisher:   NewPublisher(),
		subscribers: make(map[string][]*Subscriber),
	}
}

func (s *SubPub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub := NewSubscriber(cb)
	s.subscribers[subject] = append(s.subscribers[subject], sub)
	return sub, nil
}

func (s *SubPub) Publish(subject string, msg interface{}) error {
	return s.publisher.Publish(subject, msg)
}

func (s *SubPub) Close(ctx context.Context) error {
	for _, subs := range s.subscribers {
		for _, sub := range subs {
			sub.Unsubscribe()
		}
	}
	return nil
}
