package subpub

import (
	"SubPub_project/pkg/logger"
	"context"
	"sync"
)

type MessageHandler func(msg interface{})

type Subscriber struct {
	cb    MessageHandler
	ch    chan interface{}
	done  chan struct{}
	once  sync.Once
	close func()
}

type Subscription interface {
	Unsubscribe()
}

func NewSubscriber(ctx context.Context, cb MessageHandler) *Subscriber {
	sub := &Subscriber{
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

			case <-ctx.Done():
				sub.close()
				return

			case <-sub.done:
				return
			}
		}
	}()
	return sub
}

func (s *Subscriber) Unsubscribe() {
	logger.NewLogger().Info("Отписка от подписки:")
	s.close()
}
