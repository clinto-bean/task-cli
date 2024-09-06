package main

import (
	"log"
	"log/slog"
)

func main() {
	slog.Info("Task CLI starting. Type HELP for a list of commands.")
	db, err := NewDB("./db.json")
	if err != nil {
		log.Fatal(err)
	}
	a := API{db: db}
	a.StartREPL()
}
