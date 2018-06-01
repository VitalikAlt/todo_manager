package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lillilli/todo_manager/client/backend/config"
	"github.com/op/go-logging"
)

// Service - описание сервиса
type Service struct {
	cfg    config.Config
	mux    *mux.Router
	server *http.Server
	log    *logging.Logger
	cancel context.CancelFunc
}

// New - создание нового сервиса
func New(config config.Config) (s *Service, err error) {
	s = &Service{
		cfg: config,
		log: logging.MustGetLogger("service"),
	}

	s.mux = mux.NewRouter()
	s.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port),
		Handler:        s.mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s, nil
}

// Start - запуск нового сервиса
func (s *Service) Start() error {
	s.log.Info("Starting,...")

	ctx := context.Background()
	ctx, s.cancel = context.WithCancel(ctx)

	s.declareRoutes()
	s.log.Info("Start listen on %s", s.server.Addr)
	go s.server.ListenAndServe()

	s.log.Info("Starting... done")
	return nil
}

func (s Service) declareRoutes() {
	s.mux.PathPrefix("/").Handler(http.FileServer(http.Dir("../../../frontend/dist")))
	s.mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../../../frontend/dist"))))
}

// Stop - остановка сервиса
func (s *Service) Stop() error {
	s.log.Info("Stopping...")
	s.cancel()

	s.log.Info("Stopping... done")
	return nil
}
