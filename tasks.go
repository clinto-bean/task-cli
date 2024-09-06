package main

import (
	"fmt"
	"log"
	"log/slog"
)

type Task struct {
	Description string `json:"description"`
	Complete    bool   `json:"complete"`
}

func (api *API) CreateTask(desc string) error {
	task := Task{
		Description: desc,
	}
	t, err := api.db.AddTask(task)
	if err != nil {
		log.Printf("Could not write to DB: %v\n", err.Error())
		return nil
	}
	log.Printf("Task created: %v\n", t)
	return nil
}

func (api *API) GetTasks() {
	data, err := api.db.GetAllTasks()
	if err != nil {
		panic(err)
	}
	slog.Info("Displaying all tasks below!")
	for k, v := range data {
		msg := fmt.Sprintf("%v - %s", k, v)
		fmt.Println(msg)
	}
}

func (api *API) GetTask(id int) {
	task, err := api.db.GetTask(id)
	if err != nil {
		log.Printf("An error occurred while getting task %v: %v\n", id, err.Error())
	}
	slog.Info(fmt.Sprintf("Task %v found.", id))
	fmt.Printf("Task: %v\n", task.Description)
}

func (api *API) CompleteTask(id int) {
	err := api.db.CompleteTask(id)
	if err != nil {
		log.Printf("An error occurred while marking task %d as complete: %v\n", id, err.Error())
	}
}

func (api *API) DeleteTask(id int) {
	err := api.db.DeleteTask(id)
	if err != nil {
		log.Printf("Could not delete task %d: %v", id, err.Error())
	}
}

func (api *API) EditTask(id int, desc string) {
	err := api.db.EditTask(id, desc)
	if err != nil {
		log.Printf("Could not edit task %d: %v", id, err.Error())
	}
}
