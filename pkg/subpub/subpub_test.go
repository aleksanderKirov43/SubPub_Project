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

func TestClose(t *testing.T) {
	pubsub := NewSubPub()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pubsub.Close(ctx)
	if err != nil {
		t.Fatalf("Ошибка закрытия системы: %v", err)
	}
}
