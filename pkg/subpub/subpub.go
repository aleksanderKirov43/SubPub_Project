package subpub

import (
	"context"
	"errors"
	"log"
	"sync"
)

type MessageHandler func(msg interface{})

type Subscription interface {
	Unsubscribe()
}

type SubPub interface {
	Subscribe(subject string, cd MessageHandler) (Subscription, error)
	Publish(subject string, msg interface{}) error
	Close(ctx context.Context) error
}

func NewSubPub() SubPub {
	return &subPub{
		subjects: make(map[string][]*subscriber),
		closeCh:  make(chan struct{}),
	}
}

type subscriber struct {
	cb    MessageHandler
	ch    chan interface{}
	done  chan struct{}
	once  sync.Once
	close func()
}

type subPub struct {
	mu       sync.RWMutex
	subjects map[string][]*subscriber
	closed   bool
	closeCh  chan struct{}
}

func (s *subPub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	log.Println("Создание подписки для:", subject)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, errors.New("Подписка уже закрыта, пока нельзя добавить новые")
	}

	sub := &subscriber{
		cb:   cb,
		ch:   make(chan interface{}, 32),
		done: make(chan struct{}),
	}

	sub.close = func() {
		sub.once.Do(func() {
			close(sub.done)
		})
	}

	go func() {
		defer sub.close()
		for {
			select {
			case msg, ok := <-sub.ch:
				if !ok {
					return
				}
				sub.cb(msg)
			case <-s.closeCh:
				return
			case <-sub.done:
				return
			}
		}
	}()

	s.subjects[subject] = append(s.subjects[subject], sub)
	log.Println("Подписка зарегистрирована для:", subject)
	return &subscription{s, subject, sub}, nil
}

func (s *subPub) Publish(subject string, msg interface{}) error {
	log.Println("Публикация вызвана:", subject, msg)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return errors.New("Система закрыта")
	}

	subs := append([]*subscriber{}, s.subjects[subject]...)
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

func (s *subPub) Close(ctx context.Context) error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	close(s.closeCh)
	subs := make([]*subscriber, 0)
	for _, list := range s.subjects {
		subs = append(subs, list...)
	}
	s.subjects = nil
	s.mu.Unlock()

	done := make(chan struct{})
	go func() {
		for _, sub := range subs {
			sub.close()
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

type subscription struct {
	sp      *subPub
	subject string
	sub     *subscriber
}

func (s *subscription) Unsubscribe() {
	log.Println("Отписка от ключа:", s.subject)
	s.sp.mu.Lock()
	defer s.sp.mu.Unlock()

	list := s.sp.subjects[s.subject]
	for i, sub := range list {
		if sub == s.sub {
			s.sub.close()
			s.sp.subjects[s.subject] = append(list[:i], list[i+1:]...)
			break
		}
	}
}
