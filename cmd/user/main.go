package main

import (
	"context"
	"flag"
	"github.com/gobackpack/rmq"
	"github.com/semirm-dev/faceit/cmd/user/publisher"
	"github.com/semirm-dev/faceit/internal/db"
	"github.com/semirm-dev/faceit/user"
	"github.com/semirm-dev/faceit/user/repository"
	"github.com/sirupsen/logrus"
)

const defaultConnStr = "host=localhost port=5432 dbname=faceit_db user=postgres password=postgres sslmode=disable"

var (
	addr       = flag.String("addr", ":8001", "User Account Service address")
	rmqHost    = flag.String("rmq_host", "localhost", "RabbitMQ host address")
	connString = flag.String("connStr", defaultConnStr, "User Accounts Service connection string")
)

func main() {
	flag.Parse()

	cred := rmq.NewCredentials()
	cred.Host = *rmqHost
	hub := rmq.NewHub(cred)

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	_, err := hub.Connect(rootCtx)
	if err != nil {
		logrus.Fatal(err)
	}

	svc := user.NewAccountService(*addr, repository.NewPgDb(db.PostgresDb(*connString)), publisher.NewAccountPublisher(rootCtx, hub))

	svc.ListenForConnections(rootCtx)
}
