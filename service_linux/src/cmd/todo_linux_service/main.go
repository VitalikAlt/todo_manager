package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"os/signal"

	"github.com/lillilli/todo_manager/service_linux/src/config"
	"github.com/lillilli/todo_manager/service_linux/src/service"
	"github.com/lillilli/vconf"
)

func main() {
	var (
		configFile string
		err        error
	)

	f := flag.String("config", "", "set service config file")
	flag.Parse()

	configFile = *f

	cfg := &config.Config{}

	if err := vconf.InitFromFile(configFile, cfg); err != nil {
		fmt.Printf("unable to load config: %s\n", err)
		os.Exit(1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	service, err := service.New(cfg)
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
