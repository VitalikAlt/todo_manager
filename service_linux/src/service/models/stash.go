package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

// NewStash - создание нового хранилища
func NewStash(filePath string) *Stash {
	return &Stash{[]*Task{}, filePath}
}

// Stash - хранилище задач
type Stash struct {
	tasks    []*Task
	filePath string
}

// Set - обновить все задачи в хранилище
func (s *Stash) Set(tasks []*Task) {
	s.tasks = []*Task{}

	for _, task := range tasks {
		s.tasks = append(s.tasks, task)
	}

	s.updateFile()
}

// Add - добавление задачи в хранилище
func (s *Stash) Add(task *Task) {
	s.tasks = append(s.tasks, task)
	s.updateFile()
}

// Update - обновление задачи в хранилище
func (s *Stash) Update(task *Task) {
	if num := indexOfTask(s.tasks, task.ID); num != -1 {
		s.tasks[num] = task
		s.updateFile()
	}
}

// Reorder - обновление порядка задач в хранилище
func (s *Stash) Reorder(ids []int) {
	newTasks := []*Task{}

	for _, id := range ids {
		if num := indexOfTask(s.tasks, id); num != -1 {
			newTasks = append(newTasks, s.tasks[num])
		}
	}

	s.tasks = newTasks
	s.updateFile()
}

// Delete - удаление задачи из хранилища
func (s *Stash) Delete(id int) {
	log.Println(id)
	if num := indexOfTask(s.tasks, id); num != -1 {
		s.tasks = s.tasks[:num+copy(s.tasks[num:], s.tasks[num+1:])]
		s.updateFile()
	}
}

func (s *Stash) updateFile() {
	var content string
	log.Printf("[INFO] stash: Update file from stash, current: %#v", s.tasks)

	for _, task := range s.tasks {
		content += fmt.Sprintf("%s %v\n", task.Text, timeToString(task.DueDate))
	}

	if err := ioutil.WriteFile("/tmp/conky_notes.txt", []byte(content), 0644); err != nil {
		log.Println(err)
	}
}

func timeToString(time *time.Time) string {
	if time == nil {
		return ""
	}

	return ""
}

func indexOfTask(tasks []*Task, id int) int {
	for num, v := range tasks {
		if v.ID == id {
			return num
		}
	}

	return -1
}
