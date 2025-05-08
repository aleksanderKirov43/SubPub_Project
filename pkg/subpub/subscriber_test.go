package subpub

import "testing"

func TestUnsubscribe(t *testing.T) {
	pubsub := NewSubPub()
	sub, _ := pubsub.Subscribe("test-key", func(msg interface{}) {})

	sub.Unsubscribe()

	err := pubsub.Publish("test-key", "Hello, world!")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}
}
