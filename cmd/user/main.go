package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/faceit/user"
	"github.com/semirm-dev/faceit/user/repository"
)

var addr = flag.String("addr", ":8001", "User Account Service address")

func main() {
	flag.Parse()

	svc := user.NewAccountService(*addr, repository.NewAccountInmemory())

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	svc.ListenForConnections(rootCtx)
}
