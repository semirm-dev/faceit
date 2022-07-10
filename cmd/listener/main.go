package main

import (
	"context"
	"flag"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/faceit/event"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var rmqHost = flag.String("rmq_host", "localhost", "RabbitMQ host address")

func main() {
	flag.Parse()

	// connect
	cred := rmq.NewCredentials()
	cred.Host = *rmqHost
	hub := rmq.NewHub(cred)
	hub.ReconnectTime(30 * time.Second)

	hubCtx, hubCancel := context.WithCancel(context.Background())
	defer hubCancel()

	_, err := hub.Connect(hubCtx)
	if err != nil {
		logrus.Fatal(err)
	}

	// setup
	confAccountCreated := rmq.NewConfig()
	confAccountCreated.Exchange = "account"
	confAccountCreated.Queue = event.AccountCreated
	confAccountCreated.RoutingKey = event.AccountCreated

	if err = hub.CreateQueue(confAccountCreated); err != nil {
		logrus.Fatal(err)
	}

	confAccountModified := rmq.NewConfig()
	confAccountModified.Exchange = "account"
	confAccountModified.Queue = event.AccountModified
	confAccountModified.RoutingKey = event.AccountModified

	if err = hub.CreateQueue(confAccountModified); err != nil {
		logrus.Fatal(err)
	}

	// start consumers
	accountCreatedConsumer := hub.StartConsumer(hubCtx, confAccountCreated)
	accountModifiedConsumer := hub.StartConsumer(hubCtx, confAccountModified)

	// handle messages
	go handleConsumerMessages(hubCtx, accountCreatedConsumer, event.AccountCreated)
	go handleConsumerMessages(hubCtx, accountModifiedConsumer, event.AccountModified)

	logrus.Info("listening for messages...")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func handleConsumerMessages(ctx context.Context, cons *rmq.Consumer, name string) {
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
