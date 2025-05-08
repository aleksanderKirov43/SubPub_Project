package subpub

import (
	"context"
	"testing"
)

func TestSubscribe(t *testing.T) {
	ctx := context.Background()
	pubsub := NewSubPub()

	sub, err := pubsub.Subscribe(ctx, "test-key", func(msg interface{}) {
		t.Logf("Получено сообщение: %v", msg)
	})

	if err != nil {
		t.Fatalf("Ошибка подписки: %v", err)
	}
	if sub == nil {
		t.Fatalf("Подписка должна быть создана")
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

	if err = pubsub.Publish(ctx, "test-key", "Hello after close"); err == nil {
		t.Errorf("Ошибка! Система закрыта, но публикация всё ещё работает")
	}
}
