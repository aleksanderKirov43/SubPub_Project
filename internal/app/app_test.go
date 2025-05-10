package app

import (
	"context"
	"testing"
	"time"

	"SubPub_project/pkg/logger"
)

func TestAppPublishAndSubscribe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log := logger.NewLogger()

	appInstance := NewApp(ctx, log)

	subject := "test-subject"
	messagePublished := "hello, teamVK"

	ch := make(chan string, 1)

	subscription, err := appInstance.Subscribe(ctx, subject, func(msg interface{}) {
		if s, ok := msg.(string); ok {
			ch <- s
		}
	})
	if err != nil {
		t.Fatalf("Ошибка подписки: %v", err)
	}
	defer subscription.Unsubscribe()

	err = appInstance.Publish(ctx, subject, messagePublished)
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

	select {
	case received := <-ch:
		if received != messagePublished {
			t.Errorf("Ожидали сообщение %q, но получено %q", messagePublished, received)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Сообщение не получено подписчиком")
	}
}
