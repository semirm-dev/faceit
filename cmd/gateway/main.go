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

	accApi := gateway.NewApi(*accountAddr)

	router.POST("users", accApi.CreateAccount())
	router.PUT("users/:id", accApi.ModifyAccount())
	router.DELETE("users/:id", accApi.DeleteAccount())
	router.GET("users", accApi.GetAccounts())

	web.ServeHttp(*httpAddr, "gateway", router)
}
