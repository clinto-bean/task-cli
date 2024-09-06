package main

import (
	"bufio"
	"fmt"
	"os"
)

type API struct {
	db *DB
}

func (api *API) StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "help", "commands", "-h":
			api.CommandHelp()
		case "add", "new":
			api.CommandAdd(scanner)
		case "delete", "remove":
			api.CommandDelete(scanner)
		case "edit", "update", "change":
			api.CommandEdit(scanner)
		case "show", "tasks":
			api.CommandShow(scanner)
		case "exit", "close", "quit":
			api.CommandExit(scanner)
		case "get":
			api.CommandGetTask(scanner)
		case "complete":
			api.CommandCompleteTask(scanner)
		case "":
			fmt.Print()
		default:
			fmt.Printf("Unknown input [%v]. Type HELP for a list of commands.\n", input)
		}
	}
}
