package events

import (
	"context"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/faceit/event"
)

type accountCreated struct {
	hub *rmq.Hub
}

func NewAccountCreatedListener(hub *rmq.Hub) *accountCreated {
	return &accountCreated{
		hub: hub,
	}
}

func (ev *accountCreated) Listen(ctx context.Context) {
	consumer := startConsumer(ctx, ev.hub, event.AccountCreated)
	go handleMessages(ctx, consumer, event.AccountCreated)
}
