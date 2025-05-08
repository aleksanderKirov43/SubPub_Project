package subpub

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	ctx := context.Background()
	pubsub := NewSubPub()
	done := make(chan struct{})

	log.Println("Тест: Создаём подписку на `test-key`")
	_, _ = pubsub.Subscribe(ctx, "test-key", func(msg interface{}) {
		done <- struct{}{}
	})

	err := pubsub.Publish(ctx, "test-key", "Hello, VKTeam!")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

	select {
	case <-done:
		t.Logf("Сообщение доставлено подписчику!")
	case <-time.After(time.Second):
		t.Errorf("Таймаут! Сообщение не доставлено подписчику")
	}
}
