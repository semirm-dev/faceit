package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/faceit/event"
	"github.com/sirupsen/logrus"
)

type accountPublisher struct {
	hub  *rmq.Hub
	conf map[string]*rmq.Publisher
}

func NewAccountPublisher(ctx context.Context, hub *rmq.Hub) *accountPublisher {
	pub := &accountPublisher{
		hub:  hub,
		conf: make(map[string]*rmq.Publisher),
	}

	pub.setupQueues(ctx, []string{event.AccountCreated, event.AccountModified, event.AccountDeleted})

	return pub
}

func (pub *accountPublisher) Publish(event string, msg interface{}) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	rmqPub := pub.conf[event]
	if rmqPub == nil {
		return errors.New("invalid account event")
	}

	pub.hub.Publish(b, rmqPub)

	return nil
}

func (pub *accountPublisher) setupQueues(ctx context.Context, events []string) {
	for _, event := range events {
		conf := rmq.NewConfig()
		conf.Exchange = "account"
		conf.Queue = event
		conf.RoutingKey = event

		if err := pub.hub.CreateQueue(conf); err != nil {
			logrus.Fatal(err)
		}

		pub.conf[event] = pub.hub.CreatePublisher(ctx, conf)
	}
}
