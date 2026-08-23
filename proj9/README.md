# tri

`tri` is a command line tool. It manages tasks in a PostgreSQL database.

## Requirements

- Go 1.25 or later
- A PostgreSQL server
- Liquibase on your PATH (only for `migrateDb`)

## Setup

1. Create the file `Database/.env`. Set these variables:

   ```
   DB_HOST=<database host>
   DB_PORT=5432
   DB_USER=<database user>
   DB_PASSWORD=<database password>
   DB_NAME=todos
   ```

   Do not commit this file. The repository `.gitignore` excludes it.

2. Build the tool:

   ```
   go build -o tri .
   ```

3. Apply the database migrations:

   ```
   ./tri migrateDb
   ```

   This command creates the `todos` database if it does not exist. It then applies all Liquibase changesets. You can run this command more than one time. It is safe.

Note: the tool also creates the database on first connection if the database is missing. The `migrateDb` command is still necessary. It creates the `todos` table and its constraints.

## Commands

### Add a task

```
./tri add "Write the report"
./tri add -p 2 "Fix the build"
./tri add -d 2026-09-01 "File the report"
./tri add --due tomorrow "Call the vendor"
```

The `-p` / `--priority` flag sets the priority. Use `0` for Low, `1` for Medium, `2` for High. The default value is Low.

The `-d` / `--due` flag sets the due date. Use the format `YYYY-MM-DD`, or the words `today` or `tomorrow`. Tasks with no `--due` flag have no due date.

You can add more than one task in one command. Each argument becomes one task. All tasks in one command get the same priority and due date.

Task text must be unique. If you add a task with text that exists, the tool updates that task. It does not create a second task.

### List tasks

```
./tri list
./tri list --done
./tri list -n 7
./tri list -q report
./tri list --interactive
```

| Flag | Effect |
| --- | --- |
| `--done` | Show only completed tasks. |
| `-n`, `--near <days>` | Show only tasks with a due date in the next `<days>` days. |
| `-q`, `--query <text>` | Show only tasks that contain `<text>` in the label or text. |
| `--interactive` | Open the task list in a terminal UI. |

The list is sorted by priority. The number before each task is its position. The `done` and `delete` commands use this position.

In the interactive view:
- Use the arrow keys or `j` / `k` to move.
- Press `Enter` or `Space` to toggle the done state of the selected task.
- Press `d` to delete the selected task.
- Press `q` or `Ctrl+C` to quit.

The interactive view reloads all tasks after each change. Filters from the command line do not persist after a change.

### Complete a task

```
./tri done 1
./tri do 1
```

The argument is the position from `tri list`.

### Delete a task

```
./tri delete 1
./tri del 1
```

The argument is the position from `tri list`. If the number is larger than the list, the tool treats it as a raw database ID.

### Clear completed tasks

```
./tri clear
./tri clr
```

This command removes all tasks that are marked as done.

## Tests

```
go test ./...
```

The tests for the `Database` package require a reachable PostgreSQL server with the configuration from `Database/.env`.

## Project structure

- `cmd/` — CLI commands (Cobra).
- `Database/` — connection setup and Liquibase changelogs. See `Database/README.md`.
- `todo/` — task logic and the repository layer. See `todo/README.md`.
- `models/` — the `Todo` data structure.
- `ui/` — the Bubble Tea model for the interactive list.
- `main.go` — the entry point.
