package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/lillilli/todo_manager/service_linux/src/service/models"
	"github.com/lillilli/todo_manager/service_linux/src/service/ws/handler"
	"github.com/pkg/errors"
)

// Server - структура web socket сервера
type Server struct {
	url     url.URL
	manager *handler.Manager
	conn    *websocket.Conn
	user    *models.User

	cancel context.CancelFunc
}

// New - создание нового интсанса ws сервера
func New(addr string, user *models.User) *Server {
	url := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	return &Server{url: url, manager: handler.NewManager(), user: user}
}

// Start - начать принимать и обрабатывать ws сообщения
func (s *Server) Start() (err error) {
	log.Println("[INFO] ws: Starting ...")
	log.Printf("[INFO] ws: Connecting to %s", s.url.String())

	if s.conn, _, err = websocket.DefaultDialer.Dial(s.url.String(), nil); err != nil {
		return errors.Wrap(err, "can`t connect to ws server")
	}

	s.user.Conn = s.conn
	ctx := context.Background()
	ctx, s.cancel = context.WithCancel(ctx)

	go s.handleIncoming(ctx)
	s.registrate()

	return nil
}

func (s Server) registrate() {
	if s.user.ID != 0 {
		s.conn.WriteJSON(models.Message{Type: "auth", Data: s.user})
	} else {
		s.conn.WriteJSON(models.Message{Type: "sign_in"})
	}
}

func (s Server) handleIncoming(ctx context.Context) {
	defer s.conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg models.Message

			if err := s.conn.ReadJSON(&msg); err != nil {
				log.Printf("read json error: %v", err)
				return
			}

			log.Printf("[DEBUG] ws: Recieve message with type %s and data: %#v", msg.Type, msg.Data)

			msgContent, err := json.Marshal(msg.Data)
			if err != nil {
				log.Printf("can`t parse message data: %v", err)
				// s.conn.WriteJSON(models.Message{Type: "error", Data: []byte("can`t parse message data")})
			}

			handler := s.manager.GetHander(msg.Type)
			if err := handler.Handle(s.user, msgContent); err != nil {
				log.Println(err)
				// s.conn.WriteJSON(models.Message{Type: "error", Data: err.Error()})
			}
		}
	}
}

// Stop - завершение работы ws сервера
func (s Server) Stop() error {
	s.cancel()
	return nil
}
