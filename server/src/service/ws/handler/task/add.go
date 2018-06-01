package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/lillilli/todo_manager/server/src/db"
	"github.com/lillilli/todo_manager/server/src/db/table"
	"github.com/lillilli/todo_manager/server/src/service/ws/hub"
)

// AddHandler - обработчик task_add сообщений
type AddHandler struct {
	DB db.DB
}

// Handle - обработать сообщение
func (t AddHandler) Handle(client *hub.Client, data []byte) error {
	var params table.Task

	if ok := client.CheckAuth(); !ok {
		return client.SendAuthError()
	}

	if err := json.Unmarshal(data, &params); err != nil {
		return client.SendParseError("task_get")
	}

	fmt.Println(params.Text, params.Priority, params.DueDate)
	task, err := t.DB.Tables().Tasks.Add(client.GetUserID(), params.Text, params.Priority, params.DueDate)
	if err != nil {
		log.Printf("[ERROR] get tasks ws handler: can`t add task: %v", err)
		return errors.New("can`t add task")
	}

	return client.SendJSONbyUID("task_add", task)
}
