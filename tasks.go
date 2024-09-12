package main

import (
	"fmt"
	"time"
)

type Task struct {
	Description string `json:"description"`
	Status      string `json:"status"`
	ID          int
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Create a Task struct and attempt to add it to the database
func (api *API) CreateTask(desc string) {
	task := Task{
		Description: desc,
	}
	t, err := api.db.AddTask(task)
	if err != nil {
		api.HandleError(err)
		return
	}
	api.log.Info(fmt.Sprintf("Task created: %v", t.ID))
}

// Edit the specified task by changing its desc to match the desc passed
func (api *API) EditTask(id int, desc string) {
	err := api.db.EditTask(id, desc)
	if err != nil {
		api.HandleError(err)
	}
	api.log.Info(fmt.Sprintf("Task %d description updated: %s", id, desc))
}

// Attempts to get and print all tasks
func (api *API) GetAllTasks() {
	tasks, err := api.db.GetAllTasks()
	if err != nil {
		api.HandleError(err)
		return
	}
	for _, task := range tasks {
		message := fmt.Sprintf("\033[33m%d\033[0m %s - [%s]\n", task.ID, task.Description, task.Status)
		api.log.Printf(message)
	}
}

// Using the id, searches for a specific task
func (api *API) GetTask(id int) {
	task, err := api.db.GetTask(id)
	if err != nil {
		api.HandleError(err)
		return
	}
	api.log.Info(fmt.Sprintf("Task %v found.", id))
	api.log.Printf("Task: %v\n", task.Description)
}

// Attempts to mark a specified task as complete
func (api *API) CompleteTask(id int) {
	err := api.db.CompleteTask(id)
	if err != nil {
		api.HandleError(err)
		return
	}
	api.log.Info(fmt.Sprintf("Task %d marked as complete", id))
}

// Attempts to delete a specified task
func (api *API) DeleteTask(id int) {
	err := api.db.DeleteTask(id)
	if err != nil {
		api.HandleError(err)
		return
	}
	api.log.Info(fmt.Sprintf("Task %d successfully deleted", id))
}

// Locates and prints all tasks where Status == Complete
func (api *API) ShowCompletedTasks() {
	tasks, err := api.db.GetCompletedTasks()
	if err != nil {
		api.HandleError(err)
		return
	}
	if len(tasks) == 0 {
		api.log.Info("No tasks are complete yet!")
		return
	}
	for _, task := range tasks {
		if task.Status == "Complete" {
			message := fmt.Sprintf("\033[33m%d\033[0m %s - [%s]\n", task.ID, task.Description, task.Status)
			api.log.Printf(message)
		}
	}
}

// Locates and prints all tasks where Status != Complete
func (api *API) ShowIncompleteTasks() {
	tasks, err := api.db.GetIncompleteTasks()
	if err != nil {
		api.HandleError(err)
		return
	}
	if len(tasks) == 0 {
		api.log.Info("No incomplete tasks were found. Good job!")
		return
	}
	for _, task := range tasks {
		message := ""
		if task.Status != "Complete" {
			message = fmt.Sprintf("\033[33m%d\033[0m %s - [%s]\n", task.ID, task.Description, task.Status)
		}
		api.log.Printf(message)
	}
}

// Locates and prints all tasks where Status == In Progress
func (api *API) ShowStartedTasks() {
	tasks, err := api.db.GetIncompleteTasks()
	if err != nil {
		api.HandleError(err)
		return
	}
	if len(tasks) == 0 {
		api.log.Info("No tasks have been started yet!")
		return
	}
	for _, task := range tasks {
		message := ""
		if task.Status == "In Progress" {
			message = fmt.Sprintf("\033[33m%d\033[0m %s - [%s]\n", task.ID, task.Description, task.Status)
		}
		api.log.Printf(message)
	}
}

// Marks a task as not started/incomplete
func (api *API) UndoTask(id int) {
	err := api.db.IncompleteTask(id)
	if err != nil {
		api.HandleError(err)
		return
	}
	api.log.Info(fmt.Sprintf("Task %d successfully marked as NOT complete.", id))
}

// Marks a task as in progress
func (api *API) StartTask(id int) {
	err := api.db.StartTask(id)
	if err.Error() == "already complete" {
		api.log.Warn(fmt.Sprintf("Task %d already complete.", id))
		return
	}
	if err != nil {
		api.HandleError(err)
		return
	}
	api.log.Info(fmt.Sprintf("Task %d successfully started.", id))
}
