package service

import (
	"fmt"

	"github.com/lillilli/todo_manager/server/src/config"
	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/service/ws"
	logging "github.com/op/go-logging"
	"github.com/pkg/errors"
)

// Service - описание сервиса
type Service struct {
	db       db.DB
	cfg      *config.Config
	log      *logging.Logger
	wsServer *ws.Server
}

// New - создание нового сервиса
func New(cfg *config.Config) (service *Service, err error) {
	service = &Service{
		db:  db.NewDb(cfg.Db),
		log: logging.MustGetLogger("service"),
		cfg: cfg,
	}

	service.wsServer = ws.New(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), service.db)

	return service, nil
}

// Start - запуск нового сервиса
func (s *Service) Start() error {
	s.log.Info("Starting...")

	if err := s.db.Connect(); err != nil {
		return errors.Wrap(err, "unable to create connection to db")
	}

	if err := s.wsServer.Start(); err != nil {
		return errors.Wrap(err, "unable to start ws server")
	}

	s.log.Info("Starting... done")
	return nil
}

// Stop - остановка сервиса
func (s *Service) Stop() error {
	s.log.Info("Stopping...")

	if err := s.db.Disconnect(); err != nil {
		s.log.Errorf("Error on db disconnect: %v", err)
	}

	if err := s.wsServer.Stop(); err != nil {
		s.log.Errorf("Error on ws server close: %v", err)
	}

	s.log.Info("Stopping... done")
	return nil
}

// RunMigrations - запуск миграци БД,
//		steps 	<int> указывает количество шагов миграции к примененеию. При значении 0 применятся все доступные миграции
//		drop 	<bool> позволяет откатить все миграци
func (s *Service) RunMigrations(steps int, drop bool) (err error) {
	s.log.Info("Start migrations ...")

	if err = s.db.Connect(); err != nil {
		return errors.Wrap(err, "unable to connect to db")
	}

	if drop {
		err = s.db.Downgrade()
	} else {
		err = s.db.Migrate(steps)
	}

	if err != nil {
		return errors.Wrap(err, "migrations failed")
	}

	s.log.Info("Migrations complete")
	return nil
}
