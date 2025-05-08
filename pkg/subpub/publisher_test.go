package subpub

import (
	"testing"
	"time"
)

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

	time.Sleep(200 * time.Millisecond)

	if !received {
		t.Errorf("Сообщение не получено подписчиком")
	}
}
