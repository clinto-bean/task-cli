# Task CLI

This is a CLI written in Go used for managing to-dos. This was created as practice for CRUD apps.

## Usage

To use this CLI tool, simply clone the repository and either compile the program using Go build or alternatively run it using go run.

## Commands

`add` - Typing "Add {description}" will add a task to the database

`show` - Typing "Show all" will display all tasks, "Show complete" will show all complete tasks, "Show incomplete" will show all incomplete tasks, and "Show {id}" will show a specific task

`delete` - Typing "delete {id}" will delete the task from the db, if it exists

`complete` - Typing "complete {id}" marks task as complete but does not remove it

`undo` - Typing "undo {id}" marks a task as incomplete

`exit` - Closes the application, calling os.Exit(0) and ensuring graceful shutdown

`help` - Displays available commands

`edit` - Typing "edit {id}" allows user to edit the description

## Contribution

If for any reason you would like to contribute to this project, make a pull request.
