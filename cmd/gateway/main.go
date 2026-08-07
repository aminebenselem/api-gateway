package main

import (
	"log"
	"net/http"

	"github.com/aminebenselem/api-gateway/internal/config"
)

func main() {
	gatewayCfg, routesCfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	_ = routesCfg // we'll use this next

	server := &http.Server{
		Addr: ":" + gatewayCfg.Port,
	}

	log.Printf("Gateway listening on :%s", gatewayCfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}