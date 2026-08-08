package main

import (
	"log"
	"net/http"

	"github.com/aminebenselem/api-gateway/internal/config"
	"github.com/aminebenselem/api-gateway/internal/router"
)

func main() {
	gatewayCfg, routesCfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := router.New(routesCfg)

	server := &http.Server{
		Addr:    ":" + gatewayCfg.Port,
		Handler: r,
	}

	log.Printf("Gateway listening on :%s", gatewayCfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}