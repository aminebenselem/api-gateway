package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aminebenselem/api-gateway/internal/config"
	"github.com/aminebenselem/api-gateway/internal/middleware"
	"github.com/aminebenselem/api-gateway/internal/router"
)

func main() {
	gatewayCfg, routesCfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := router.New(routesCfg)

	handler := middleware.Logging(r)

	server := &http.Server{
		Addr:    ":" + gatewayCfg.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("Gateway listening on :%s", gatewayCfg.Port)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutting down gateway...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Gateway shutdown error: %v", err)
	}

	log.Println("Gateway stopped")
}