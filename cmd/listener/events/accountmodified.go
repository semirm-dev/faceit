package events

import (
	"context"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/faceit/event"
)

type accountModified struct {
	hub *rmq.Hub
}

func NewAccountModifiedListener(hub *rmq.Hub) *accountModified {
	return &accountModified{
		hub: hub,
	}
}

func (ev *accountModified) Listen(ctx context.Context) {
	consumer := startConsumer(ctx, ev.hub, event.AccountModified)
	go handleMessages(ctx, consumer, event.AccountModified)
}
