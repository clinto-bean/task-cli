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

func (l *Slogger) Error(err error) {
	fmt.Println(Red+"Error"+Reset+":", err.Error())
}

func (l *Slogger) Info(msg string) {
	msg = Green + "Info: " + Reset + msg
	fmt.Println(msg)

}

func (l *Slogger) Warn(msg string) {
	msg = Yellow + "Warning: " + Reset + msg
	fmt.Println(msg)
}

func (l *Slogger) Println(msg string) {
	fmt.Println(msg)
}

func (l *Slogger) Announce(msg string) {
	msg = Blue + msg + Reset
	fmt.Println(msg)
}

func (l *Slogger) Fatal(err error) {
	l.Error(err)
	os.Exit(1)
}

func (l *Slogger) Printf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Print(msg)
}
