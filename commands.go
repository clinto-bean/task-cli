package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Displays all tasks in a map of "task name": "task description"
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
		"start":    "marks a task as in progress. format: \"start {id}\"",
	}

	for k, v := range commands {
		api.log.Println("\033[96m" + k + "\033[0m: " + v)
	}
}

// Takes a description string and calls CreateTask passing it the description
func (api *API) CommandAdd(desc string) {
	api.CreateTask(desc)
}

// Takes an integer of task.ID and calls DeleteTask with it
func (api *API) CommandDelete(taskID int) {
	api.DeleteTask(taskID)
}

// Takes in the id of a task and the new description, passing both to Edit Task
func (api *API) CommandEdit(id int, desc string) {
	api.EditTask(id, desc)
}

// runs exec.Command("clear") then exits. Pass it any messages you wish to have printed after the program exits, such as errors that persist beyond sessions
func (api *API) CommandExit(args ...any) {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
	if len(args) > 0 {
		for _, arg := range args {
			api.log.Info(fmt.Sprintf("Important: %v", arg))
		}
	}
	os.Exit(0)
}

// Calls GetAllTasks
func (api *API) CommandShowAll() {
	api.GetAllTasks()
}

// Takes in the id of a task and calls GetTask
func (api *API) CommandGet(id int) {
	api.GetTask(id)
}

// Takes in the id of a task and calls CompleteTask
func (api *API) CommandComplete(id int) {
	api.CompleteTask(id)
}

// Calls ShowCompletedTasks
func (api *API) CommandShowComplete() {
	api.ShowCompletedTasks()
}

// Takes in the id of a task and calls UndoTask
func (api *API) CommandUndo(id int) {
	api.UndoTask(id)
}

// Calls ShowIncompleteTask
func (api *API) CommandShowIncomplete() {
	api.ShowIncompleteTasks()
}

// Takes in the id of a task and calls StartTask
func (api *API) CommandStart(id int) {
	api.StartTask(id)
}

// Calls ShowStartedTasks
func (api *API) CommandShowStarted() {
	api.ShowStartedTasks()
}
