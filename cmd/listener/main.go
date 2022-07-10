package main

import (
	"context"
	"flag"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/faceit/cmd/listener/account"
	"github.com/semirm-dev/faceit/event"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var rmqHost = flag.String("rmq_host", "localhost", "RabbitMQ host address")

func main() {
	flag.Parse()

	cred := rmq.NewCredentials()
	cred.Host = *rmqHost
	hub := rmq.NewHub(cred)

	hubCtx, hubCancel := context.WithCancel(context.Background())
	defer hubCancel()

	_, err := hub.Connect(hubCtx)
	if err != nil {
		logrus.Fatal(err)
	}

	// create listeners for different account actions/events
	accountCreatedConsumer := account.StartConsumer(hubCtx, hub, event.AccountCreated)
	accountModifiedConsumer := account.StartConsumer(hubCtx, hub, event.AccountModified)
	accountDeletedConsumer := account.StartConsumer(hubCtx, hub, event.AccountDeleted)

	// handle messages
	go account.HandleMessages(hubCtx, accountCreatedConsumer, event.AccountCreated)
	go account.HandleMessages(hubCtx, accountModifiedConsumer, event.AccountModified)
	go account.HandleMessages(hubCtx, accountDeletedConsumer, event.AccountDeleted)

	logrus.Info("listening for messages...")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
