package task

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/db/table"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// DeleteHandler - обработчик task_add сообщений
type DeleteHandler struct {
	DB db.DB
}

// Handle - обработать сообщение
func (t DeleteHandler) Handle(client *hub.Client, data []byte) error {
	var params table.Task

	if ok := client.CheckAuth(); !ok {
		return client.SendAuthError()
	}

	if err := json.Unmarshal(data, &params); err != nil {
		return client.SendParseError("task_get")
	}

	if err := t.DB.Tables().Tasks.Delete(params.ID); err != nil {
		log.Printf("[ERROR] get tasks ws handler: can`t delete task: %v", err)
		return errors.New("can`t delete task")
	}

	return client.SendJSONbyUID("task_delete", params)
}
