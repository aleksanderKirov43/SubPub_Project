package subpub

import (
	"context"
	"testing"
)

func TestSubscribe(t *testing.T) {
	pubsub := NewSubPub()
	sub, err := pubsub.Subscribe("test-key", func(msg interface{}) {
		t.Logf("Получено сообщение: %v", msg)
	})

	if err != nil {
		t.Fatalf("Ошибка подписки: %v", err)
	}
	if sub == nil {
		t.Fatalf("Подписка должна быть создана")
	}
}

func TestPublish(t *testing.T) {
	pubsub := NewSubPub()
	received := false

	_, _ = pubsub.Subscribe("test-key", func(msg interface{}) {
		received = true
	})

	err := pubsub.Publish("test-key", "Hello, world!")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

	if !received {
		t.Errorf("Сообщение не получено подписчиком")
	}
}

func TestUnsubscribe(t *testing.T) {
	pubsub := NewSubPub()
	sub, _ := pubsub.Subscribe("test-key", func(msg interface{}) {})

	sub.Unsubscribe()

	err := pubsub.Publish("test-key", "Hello, world!")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}
}

func TestClose(t *testing.T) {
	pubsub := NewSubPub()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pubsub.Close(ctx)
	if err != nil {
		t.Fatalf("Ошибка закрытия системы: %v", err)
	}
}
