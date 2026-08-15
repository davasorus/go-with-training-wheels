# Command Layer

This directory contains the command-line interface (CLI) logic. It uses the Cobra library to define and manage commands.

## Components

- **root.go**: This file defines the root command for the application. It handles global flags and initializes the database connection before executing subcommands.
- **add.go**: This file implements the `add` command. It allows users to add new items to the list with a specified priority level.
- **done.go**: This file implements the `done` command (aliased as `do`). It marks an item as complete based on its index in the list.
- **list.go**: This file implements the `list` command. It displays the current tasks in a formatted table. It can filter for only completed items.
- **migrateDb.go**: This file handles database migration logic (Note: currently placeholder/not fully utilized).

## Logic Flow

1. The user invokes a command via the terminal.
2. `root.go` initializes the application and connects to the database.
3. The specific subcommand (add, done, or list) processes the input.
4. The command interacts with the `todo` package for business logic.
5. The result is printed to the standard output.
