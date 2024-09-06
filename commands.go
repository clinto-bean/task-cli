package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
)

func (api *API) CommandHelp() {
	fmt.Println("A list of commands is found below.")
	fmt.Println("\tHelp - Show available commands")
	fmt.Println("\tAdd - create a task")
	fmt.Println("\tEdit - edit an existing task")
	fmt.Println("\tShow - shows all tasks")
	fmt.Println("\tGet - gets a specific task by its ID")
	fmt.Println("\tComplete - completes a task")
	fmt.Println("\tDelete - deletes a specific task by its ID")
	fmt.Println("\tExit - closes the program by calling os.Exit(0)")
}

func (api *API) CommandAdd(scanner *bufio.Scanner) {
	fmt.Printf("Please enter the description for your task:\n> ")
	scanner.Scan()
	description := scanner.Text()
	err := api.CreateTask(description)
	if err != nil {
		panic(err)
	}

}

func (api *API) CommandDelete(scanner *bufio.Scanner) {
	fmt.Printf("Please enter the ID of the task you wish to delete.\nIf you are unsure of the ID, type 'show' to display tasks.\n> ")
	scanner.Scan()
	input := scanner.Text()
	if input == "show" {
		api.CommandShow(scanner)
		return
	}
	taskID, err := strconv.Atoi(input)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Deleting task [%v]\n", taskID)
	api.DeleteTask(taskID)
}

func (api *API) CommandEdit(scanner *bufio.Scanner) {}

func (api *API) CommandExit(scanner *bufio.Scanner) {
	log.Println("Application closing. Goodbye!")
	os.Exit(0)
}

func (api *API) CommandShow(scanner *bufio.Scanner) {
	fmt.Println("Showing commands")
	api.GetTasks()
}

func (api *API) CommandGetTask(scanner *bufio.Scanner) {
	fmt.Printf("Please enter the ID of the task you are searching for...\n> ")
	scanner.Scan()
	input := scanner.Text()
	id, err := strconv.Atoi(input)
	if err != nil {
		slog.Error("Please enter a numeric ID")
		return
	}
	api.GetTask(id)
}

func (api *API) CommandCompleteTask(scanner *bufio.Scanner) {
	fmt.Printf("Please enter the ID of the task you are completing...\n>")
	scanner.Scan()
	input := scanner.Text()
	id, err := strconv.Atoi(input)
	if err != nil {
		slog.Error("Please enter a numeric ID")
	}
	api.CompleteTask(id)
}
