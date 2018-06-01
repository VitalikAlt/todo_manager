package task

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// ReorderHandler - обработчик task_reorder сообщений
type ReorderHandler struct {
	DB db.DB
}

// Handle - обработать сообщение
func (t ReorderHandler) Handle(client *hub.Client, data []byte) error {
	var params []int

	if ok := client.CheckAuth(); !ok {
		return client.SendAuthError()
	}

	if err := json.Unmarshal(data, &params); err != nil {
		return client.SendParseError("task_get")
	}

	if err := t.DB.Tables().Tasks.Reorder(params); err != nil {
		log.Printf("[ERROR] get tasks ws handler: can`t reorder task: %v", err)
		return errors.New("can`t reorder task")
	}

	return client.SendJSONbyUID("task_reorder", params)
}
