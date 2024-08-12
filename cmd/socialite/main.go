package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"socialite/config"
	"socialite/server"
)

func main() {
	globalCtx, cancelGlobalCtx := context.WithCancel(context.Background())

	cfg := config.ParseConfig()
	config.Verify(cfg)

	// create new server instance and start
	serverInstance := server.New(globalCtx, cfg)
	go func() {
		err := serverInstance.Start()
		if err != nil {
			log.Fatalf("[ERROR] starting server for %s : %s", cfg.Server.ServiceName, err.Error())
		}
	}()

	// wait for signal interrupt
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Printf("[INFO] stopping server for %s", cfg.Server.ServiceName)
	cancelGlobalCtx()
}
