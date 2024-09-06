# Task CLI

This is a CLI written in Go used for managing to-dos. This was created as practice for CRUD apps.

## Usage

`add` - Prompts user for a description of a task, then adds it to the database
`show` - Shows all tasks in the database in the format of `ID: Description`
`get` - Prompts user for task ID and fetches the specified task
`delete` - Prompts user for ID of a task, then removes it from the database
`complete` - Prompts user for ID of a task, then marks it as complete but does not remove it
`exit` - Closes the application, calling os.Exit(0) and ensuring graceful shutdown
`help` - Displays available commands
`edit` - Prompts user for ID of a task, then allows them to edit the description of it

License: MIT
