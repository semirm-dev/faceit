package events

import (
	"context"
	"github.com/gobackpack/rmq"
	"github.com/sirupsen/logrus"
)

func startConsumer(ctx context.Context, hub *rmq.Hub, event string) *rmq.Consumer {
	conf := rmq.NewConfig()
	conf.Exchange = "account"
	conf.Queue = event
	conf.RoutingKey = event

	if err := hub.CreateQueue(conf); err != nil {
		logrus.Fatal(err)
	}

	consumer := hub.StartConsumer(ctx, conf)

	return consumer
}

func handleMessages(ctx context.Context, cons *rmq.Consumer, name string) {
	logrus.Infof("%s started", name)

	defer logrus.Warnf("%s closed", name)

	for {
		select {
		case msg := <-cons.OnMessage:
			logrus.Infof("%s - %s", name, msg)
		case err := <-cons.OnError:
			logrus.Error(err)
		case <-ctx.Done():
			return
		}
	}
}
