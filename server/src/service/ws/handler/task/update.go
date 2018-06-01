package task

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/db/table"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// UpdateHandler - обработчик task_add сообщений
type UpdateHandler struct {
	DB db.DB
}

// Handle - обработать сообщение
func (t UpdateHandler) Handle(client *hub.Client, data []byte) error {
	var params table.Task

	if ok := client.CheckAuth(); !ok {
		return client.SendAuthError()
	}

	if err := json.Unmarshal(data, &params); err != nil {
		return client.SendParseError("task_update")
	}

	if err := t.DB.Tables().Tasks.Update(params.ID, params.Text, params.Priority, params.DueDate); err != nil {
		log.Printf("[ERROR] get tasks ws handler: can`t update task: %v", err)
		return errors.New("can`t update task")
	}

	return client.SendJSONbyUID("task_update", params)
}
