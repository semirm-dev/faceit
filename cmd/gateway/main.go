package main

import (
	"flag"
	"github.com/semirm-dev/faceit/gateway"
	"github.com/semirm-dev/faceit/internal/web"
)

var (
	httpAddr    = flag.String("http", ":8000", "Http address")
	accountAddr = flag.String("account_uri", ":8001", "User Account Service address")
)

func main() {
	flag.Parse()

	router := web.NewRouter()

	api := gateway.NewApi(*accountAddr)

	router.POST("users", api.CreateAccount())
	router.PUT("users/:id", api.ModifyAccount())
	router.DELETE("users/:id", api.DeleteAccount())
	router.GET("users", api.GetAccounts())

	web.ServeHttp(*httpAddr, "gateway", router)
}
