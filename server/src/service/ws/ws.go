package ws

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/service/ws/handler"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"polling"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Server - структура web socket сервера
type Server struct {
	hub     *hub.Hub
	db      db.DB
	addr    string
	manager *handler.Manager

	ctx    context.Context
	cancel context.CancelFunc
}

// New - создание нового интсанса ws сервера
func New(addr string, db db.DB) *Server {
	return &Server{addr: addr, db: db, hub: hub.NewHub(), manager: handler.NewManager(db)}
}

// Start - начать принимать и обрабатывать ws сообщения
func (s *Server) Start() error {
	log.Println("[INFO] ws: Starting ...")
	ctx := context.Background()
	s.ctx, s.cancel = context.WithCancel(ctx)

	http.HandleFunc("/", s.handleWS)
	go http.ListenAndServe(s.addr, nil)
	log.Printf("[INFO] ws: Start listen on ws://%s/", s.addr)

	return nil
}

// Message - протокол общения по ws
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (s Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer conn.Close()
	client := s.hub.NewClient(conn)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			var msg Message

			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("read json error: %v", err)
				return
			}

			log.Printf("[DEBUG] ws: Recieve message with type %s and data: %#v", msg.Type, msg.Data)

			msgContent, err := json.Marshal(msg.Data)
			if err != nil {
				conn.WriteJSON(Message{Type: "error", Data: []byte("can`t parse message data")})
			}

			handler := s.manager.GetHander(msg.Type)
			if err := handler.Handle(client, msgContent); err != nil {
				log.Println(err)
				conn.WriteJSON(Message{Type: "error", Data: err.Error()})
			}
		}
	}
}

// Stop - завершение работы ws сервера
func (s Server) Stop() error {
	s.cancel()
	return nil
}

func getBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
