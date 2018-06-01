package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"os/signal"

	"github.com/lillilli/todo_manager/server/src/config"
	"github.com/lillilli/todo_manager/server/src/service"
	"github.com/lillilli/vconf"
)

var (
	configFile = flag.String("config", "", "set service config file")
	migrate    = flag.Bool("migrate", false, "set migrate action")
	steps      = flag.Int("steps", 0, "steps number for migrations")
	drop       = flag.Bool("drop", false, "set drop all migrations")
)

func main() {
	flag.Parse()

	cfg := &config.Config{}

	if err := vconf.InitFromFile(*configFile, cfg); err != nil {
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

	if *migrate {
		err := service.RunMigrations(*steps, *drop)

		if err != nil {
			fmt.Printf("unable to migrate DB: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err = service.Start(); err != nil {
		fmt.Printf("unable to start service: %s\n", err)
		os.Exit(1)
	}

	<-signals
	close(signals)

	service.Stop()
}
