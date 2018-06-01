package main

import (
	"fmt"
	"os"
	"syscall"

	"os/signal"

	"github.com/lillilli/todo_manager/client/backend/config"
	"github.com/lillilli/todo_manager/client/backend/service"
	"github.com/lillilli/vconf"
)

func main() {
	cfg := &config.Config{}

	if err := vconf.Init(cfg); err != nil {
		fmt.Printf("unable to load config: %s\n", err)
		os.Exit(1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	service, err := service.New(*cfg)
	if err != nil {
		fmt.Printf("unable to create service: %s\n", err)
		os.Exit(1)
	}

	if err = service.Start(); err != nil {
		fmt.Printf("unable to start service: %s\n", err)
		os.Exit(1)
	}

	<-signals
	close(signals)

	service.Stop()
}
