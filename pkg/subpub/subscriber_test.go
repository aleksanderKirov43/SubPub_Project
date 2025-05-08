package subpub

import (
	"context"
	"testing"
)

func TestUnsubscribe(t *testing.T) {
	ctx := context.Background()
	pubsub := NewSubPub()
	received := false

	sub, _ := pubsub.Subscribe(ctx, "test-key", func(msg interface{}) {
		received = true
	})

	sub.Unsubscribe()

	err := pubsub.Publish(ctx, "test-key", "Hello, world!")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

	if received {
		t.Errorf("Ошибка! Отписанный подписчик получил сообщение")
	}
}
