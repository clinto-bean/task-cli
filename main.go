package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
	log := &Slogger{}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	log.Println("")
	log.Announce("Task CLI starting.")
	log.Info("Connecting to database.")
	db, err := NewDB("./db.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Announce("Database connected successfully.")
	log.Info("Type HELP for a list of commands.")
	a := API{db: db, log: log}
	go func() {
		for range sig {
			fmt.Println()
			log.Warn("Application interrupted. Ensuring graceful shutdown.")
			a.CommandExit()
		}
	}()
	log.Fatal(a.StartREPL())
}
