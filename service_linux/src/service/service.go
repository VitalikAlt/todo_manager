package service

import (
	"log"

	"github.com/lillilli/todo_manager/service_linux/src/config"
	"github.com/lillilli/todo_manager/service_linux/src/service/models"
	"github.com/lillilli/todo_manager/service_linux/src/service/ws"
	logging "github.com/op/go-logging"
	"github.com/pkg/errors"
)

// Service - описание сервиса
type Service struct {
	cfg      *config.Config
	log      *logging.Logger
	wsServer *ws.Server
}

// New - создание нового сервиса
func New(cfg *config.Config) (service *Service, err error) {
	service = &Service{
		log: logging.MustGetLogger("service"),
		cfg: cfg,
	}

	log.Println(cfg)

	user := &models.User{
		ID:    cfg.Client.ID,
		Token: cfg.Client.Token,
		Stash: models.NewStash("./test"),
	}

	service.wsServer = ws.New(cfg.Server.Address, user)

	return service, nil
}

// Start - запуск нового сервиса
func (s *Service) Start() error {
	s.log.Info("Starting...")

	if err := s.wsServer.Start(); err != nil {
		return errors.Wrap(err, "unable to start ws server")
	}

	s.log.Info("Starting... done")
	return nil
}

// Stop - остановка сервиса
func (s *Service) Stop() error {
	s.log.Info("Stopping...")

	if err := s.wsServer.Stop(); err != nil {
		s.log.Errorf("Error on ws server close: %v", err)
	}

	s.log.Info("Stopping... done")
	return nil
}
