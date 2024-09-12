package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// Upon startup the program clears the command-line to reduce noise
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
	// Create logger structure for use within program
	log := &Slogger{}
	// Watch for ^C and handle it gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	log.Announce("\nTask CLI starting.")
	log.Info("Connecting to database.")
	// Attempt to establish DB connection
	db, err := NewDB("./db.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Announce("Database connected successfully.")
	log.Info("Type HELP for a list of commands.")
	// Create API
	a := API{db: db, log: log}
	// Handle any ^C call
	go func() {
		for range sig {
			fmt.Println()
			log.Warn("Application interrupted. Ensuring graceful shutdown.")
			a.CommandExit()
		}
	}()
	// Start the program loop
	a.StartREPL()
}
