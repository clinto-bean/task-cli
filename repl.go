package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type API struct {
	db  *DB
	log *Slogger
}

func (api *API) StartREPL() error {

	NaN := errors.New("enter a valid numeric ID")
	BadArgs := errors.New("insufficient arguments for command. type \"help\" for info")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		args := strings.Split(input, " ")
		args[0] = strings.ToLower(args[0])

		switch args[0] {
		case "help", "commands":
			api.CommandHelp()
		case "add", "new":
			desc, _ := parseArgs(args[1:])
			api.CommandAdd(desc)
		case "delete", "remove":
			id, err := strconv.Atoi(args[1])
			if err != nil {
				api.HandleError(NaN)
				continue
			}
			api.CommandDelete(id)
			continue
		case "edit", "update", "change":
			if len(args) > 2 {
				id, err := strconv.Atoi(args[1])
				if err != nil {
					api.HandleError(err)
					continue
				}
				desc, _ := parseArgs(args[2:])
				api.CommandEdit(id, desc)
				continue
			}
			api.HandleError(BadArgs)
		case "show", "tasks":
			if len(args) > 1 {
				if args[1] == "complete" {
					api.CommandShowComplete()
					continue
				}
				if args[1] == "all" {
					api.CommandShowAll()
					continue
				}
				if args[1] == "incomplete" {
					api.CommandShowIncomplete()
					continue
				}
				if args[1] == "started" {
					api.CommandShowStarted()
					api.log.Info("started tasks")
					continue
				}
				id, err := strconv.Atoi(args[1])
				if err != nil {
					api.HandleError(NaN)
					continue
				}
				api.CommandGet(id)

			} else {
				api.log.Warn("Incorrect format. Please type 'show (\"complete\", \"all\" or \"{taskID\"}). Type 'help' for more information.")
			}
		case "exit", "close", "quit":
			api.CommandExit(nil)
		case "start":
			id, err := strconv.Atoi(args[1])
			if err != nil {
				api.HandleError(NaN)
				continue
			}
			api.CommandStart(id)
		case "complete":
			if len(args) == 1 {
				api.log.Warn("Incorrect format. Please type 'complete {id}' or 'help' for more information.")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				api.HandleError(NaN)
				continue
			}
			api.CommandComplete(id)
		case "undo":
			id, err := strconv.Atoi(args[1])
			if err != nil {
				api.HandleError(err)
				continue
			}
			api.CommandUndo(id)
		case "":
			fmt.Print()
		default:
			api.log.Printf("Unknown input [%v]. Type HELP for a list of commands.\n", input)
		}
	}
}

func parseArgs(args []string) (string, []string) {
	var words = make([]string, 0, len(args))
	var flags = make([]string, 0, len(args))
	for _, arg := range args {
		if len(arg) > 1 && arg[0:2] == "--" {
			flags = append(flags, arg)
			continue
		}
		words = append(words, arg)
	}
	msg := strings.Join(words, " ")
	return msg, flags
}

func (api *API) HandleError(err error) {
	if strings.ToLower(err.Error()[0:4]) == "fatal" {
		api.log.Fatal(err)
		return
	}
	switch err.Error() {
	case "not found":
		api.log.Warn("Task not found!")
	case "already complete":
		api.log.Warn("Task already complete!")
	default:
		api.log.Error(err)
	}
}
