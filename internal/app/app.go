package app

import (
	"SubPub_project/pkg/logger"
	"SubPub_project/pkg/subpub"

	"context"
)

type App struct {
	PubSub    subpub.SubPubInterface
	Publisher subpub.PublisherInterface
	Logger    logger.Logger
}

func NewApp(ctx context.Context, logInstance logger.Logger) *App {

	publisher := subpub.NewPublisher(ctx, logInstance)
	pubsub := subpub.NewSubPub()

	return &App{
		PubSub:    pubsub,
		Publisher: publisher,
		Logger:    logInstance,
	}
}

func (a *App) Subscribe(ctx context.Context, subject string, handler subpub.MessageHandler) (subpub.Subscription, error) {
	return a.PubSub.Subscribe(ctx, subject, handler)
}

func (a *App) Publish(ctx context.Context, subject string, msg interface{}) error {
	return a.PubSub.Publish(ctx, subject, msg)
}
