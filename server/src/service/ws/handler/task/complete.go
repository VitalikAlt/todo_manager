package task

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/db/table"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// CompleteHandler - обработчик task_complete сообщений
type CompleteHandler struct {
	DB db.DB
}

// Handle - обработать сообщение
func (t CompleteHandler) Handle(client *hub.Client, data []byte) error {
	var params table.Task

	if ok := client.CheckAuth(); !ok {
		return client.SendAuthError()
	}

	if err := json.Unmarshal(data, &params); err != nil {
		return client.SendParseError("task_get")
	}

	if err := t.DB.Tables().Tasks.Complete(params.ID); err != nil {
		log.Printf("[ERROR] get tasks ws handler: can`t complete task: %v", err)
		return errors.New("can`t complete task")
	}

	return client.SendJSONbyUID("task_complete", params)
}
