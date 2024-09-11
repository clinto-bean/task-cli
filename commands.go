package main

import (
	"os"
	"os/exec"
)

func (api *API) CommandHelp() {

	commands := map[string]string{
		"help":     "show available commands",
		"add":      "add a task. format: \"add {description}\"",
		"edit":     "edit an existing task. format: \"edit {id} {description}\"",
		"complete": "complete an existing task. format: \"complete {id}\"",
		"undo":     "marks a task as incomplete. format: \"undo {id}\"",
		"show":     "shows all tasks, \"show complete\" to display complete, \"show incomplete\" to show incomplete, or \"show {id}\" for a specific task",
		"delete":   "deletes a specific task. format: \"delete {id}\"",
		"exit":     "closes the program by calling os.Exit(0)",
	}

	for k, v := range commands {
		api.log.Println("\033[96m" + k + "\033[0m: " + v)
	}
}

func (api *API) CommandAdd(desc string) {
	api.CreateTask(desc)
}

func (api *API) CommandDelete(taskID int) {
	api.DeleteTask(taskID)
}

func (api *API) CommandEdit(id int, desc string) {
	api.EditTask(id, desc)
}

func (api *API) CommandExit(args ...any) {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
	os.Exit(0)
}

func (api *API) CommandShowAll() {
	api.GetAllTasks()
}

func (api *API) CommandGetTask(id int) {
	api.GetTask(id)
}

func (api *API) CommandComplete(id int) {
	api.CompleteTask(id)
}

func (api *API) CommandShowComplete() {
	api.ShowCompletedTasks()
}

func (api *API) CommandUndo(id int) {
	api.UndoTask(id)
}

func (api *API) CommandShowIncomplete() {
	api.ShowIncompleteTasks()
}
