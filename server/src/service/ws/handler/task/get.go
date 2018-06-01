package task

import (
	"log"

	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// GetHandler - обработчик task_get сообщений
type GetHandler struct {
	DB db.DB
}

// Handle - обработать сообщение
func (t GetHandler) Handle(client *hub.Client, data []byte) error {
	if ok := client.CheckAuth(); !ok {
		return client.SendAuthError()
	}

	tasks, err := t.DB.Tables().Tasks.GetAll(client.GetUserID())
	if err != nil {
		log.Printf("[ERROR] get tasks ws handler: can`t get tasks: %v", err)
	}

	return client.SendJSON("task_get", tasks)
}
