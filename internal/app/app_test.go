package app

import (
	"context"
	"testing"
	"time"

	"SubPub_project/pkg/subpub"
	pb "SubPub_project/proto"
)

type MockSubscriber struct {
	ch   chan interface{}
	done chan struct{}
}

func NewMockSubscriber() *MockSubscriber {
	return &MockSubscriber{
		ch:   make(chan interface{}, 1), // Буферизованный канал
		done: make(chan struct{}),
	}
}

func (m *MockSubscriber) Handler(msg interface{}) {
	m.ch <- msg
}

func (m *MockSubscriber) Close() {
	close(m.done)
}

func TestPublishRPC(t *testing.T) {
	ctx := context.Background()
	pubsub := subpub.NewSubPub()
	server := NewServer(ctx, pubsub)

	received := false

	_, _ = pubsub.Subscribe(ctx, "test-key", func(msg interface{}) {
		received = true
	})

	req := &pb.PublishRequest{
		Key:  "test-key",
		Data: "test-Data",
	}

	_, err := server.Publish(context.Background(), req)
	if err != nil {
		t.Fatalf("Ошибка публикации через gRPC: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if !received {
		t.Errorf("Сообщение не получено подписчиком")
	}
}

func TestSubscribeRPC(t *testing.T) {
	ctx := context.Background()
	pubsub := subpub.NewSubPub()

	sub := NewMockSubscriber()
	_, _ = pubsub.Subscribe(ctx, "test-key", sub.Handler) // Подключаем подписчика

	data := "test text"
	_ = pubsub.Publish(ctx, "test-key", data)

	time.Sleep(200 * time.Millisecond)

	select {
	case received := <-sub.ch:
		if received != data {
			t.Errorf("Ожидали сообщение %s, но получили %v", data, received)
		}
	default:
		t.Errorf("Сообщение не получено подписчиком")
	}
}
