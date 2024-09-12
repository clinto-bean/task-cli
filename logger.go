package main

import (
	"fmt"
	"os"
)

type Slogger struct {
}

const (
	Reset  = "\033[0;0m"
	Red    = "\033[91m"
	Green  = "\033[92m"
	Yellow = "\033[93m"
	Blue   = "\033[94m"
)

// Error adds red text to the error label
func (l *Slogger) Error(err error) {
	fmt.Println(Red+"Error"+Reset+":", err.Error())
}

// Info adds green text to the info label
func (l *Slogger) Info(msg string) {
	msg = Green + "Info: " + Reset + msg
	fmt.Println(msg)

}

// Warn adds yellow text to the warning label
func (l *Slogger) Warn(msg string) {
	msg = Yellow + "Warning: " + Reset + msg
	fmt.Println(msg)
}

// Calls to fmt.Println with the message provided
func (l *Slogger) Println(msg string) {
	fmt.Println(msg)
}

// Announce prints all text in blue
func (l *Slogger) Announce(msg string) {
	msg = Blue + msg + Reset
	fmt.Println(msg)
}

// Fatal calls l.Error then os.Exit(1)
func (l *Slogger) Fatal(err error) {
	l.Error(err)
	os.Exit(1)
}

// Prints a formatted message using fmt.Sprintf and fmt.Print
func (l *Slogger) Printf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Print(msg)
}
