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

   Note: changeset 004 adds the `created_at`, `completed_at`, and `tags` columns. Run `migrateDb` again after an update to apply new changesets.

Note: the tool also creates the database on first connection if the database is missing. The `migrateDb` command is still necessary. It creates the `todos` table and its constraints.

## Commands

### Add a task

```
./tri add "Write the report"
./tri add -p high "Fix the build"
./tri add -d 2026-09-01 "File the report"
./tri add --due tomorrow --tag home "Call the vendor"
```

The `-p` / `--priority` flag sets the priority. Use `low`, `medium`, or `high`. The numbers `0`, `1`, and `2` are also valid. The default value is `low`.

The `-d` / `--due` flag sets the due date. Use the format `YYYY-MM-DD`, or the words `today` or `tomorrow`. Tasks with no `--due` flag have no due date.

The `--tag` flag adds a tag. Repeat the flag to add more tags.

You can add more than one task in one command. Each argument becomes one task. All tasks in one command get the same priority, due date, and tags.

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
| `--tag <tag>` | Show only tasks with this tag. Repeat the flag to require more tags. |
| `-i`, `--interactive` | Open the task list in a terminal UI. |

A due date in the past shows in red with an "(overdue)" mark. The list is sorted by priority. The number before each task is its row number in the full sorted list. The `done` and `delete` commands use this number.

Caution: the `done` and `delete` commands always count against the full list. Run `tri list` with no filter flags to see the correct numbers before you use `done` or `delete`.

The interactive view opens in a full screen and restores your terminal on exit. Long lists scroll with the cursor.

In the interactive view:
- Use the arrow keys or `j` / `k` to move.
- Press `Enter` or `Space` to toggle the done state of the selected task.
- Press `a` to add a new task. Type the text, then press `Enter` to save or `Esc` to cancel.
- Press `e` to edit the text of the selected task.
- Press `1`, `2`, or `3` to set the priority of the selected task to Low, Medium, or High.
- Press `d` to delete the selected task.
- Press `q` or `Ctrl+C` to quit.

After a change, the cursor stays on the same task, even when the list order changes.

The interactive view reloads the tasks after each change. Filters from the command line stay in effect.

### Complete a task

```
./tri done 1
./tri do 1
```

The argument is the position from `tri list`.

### Mark a task as not done

```
./tri undo 1
```

The argument is the number from `tri list`. This clears the completion time.

### Edit a task

```
./tri edit 2 --text "New description"
./tri edit 2 -p high --due 2026-09-01
./tri edit 2 --due none --tags none
```

The first argument is the number from `tri list`. Only the flags you provide change. Use `--due none` to remove the due date. Use `--tags none` to remove all tags. The `--tags` flag replaces the full tag list.

### Move a task

```
./tri move 3 up
./tri move 2 down
```

Move a task one place up or down in the list. Movement is only possible within the same priority and done state. To move a task further, change its priority with `tri edit`.

### Export tasks as JSON

```
./tri export
./tri export | jq '.[] | select(.done)'
```

Writes all tasks as a JSON array to standard output. Each task includes its creation time and, for done tasks, its completion time.

### Delete a task

```
./tri delete 1
./tri del 1 --yes
```

The argument is the position from `tri list`. If the number is larger than the list, the tool treats it as a raw database ID. The command shows the task text and asks for confirmation. Use `-y` / `--yes` to skip the question, for example in scripts.

### Clear completed tasks

```
./tri clear
./tri clr --yes
```

This command removes all tasks that are marked as done. It shows the count and asks for confirmation. Use `-y` / `--yes` to skip the question.

## Exit codes and scripting

Commands return exit code `0` on success and `1` on failure. Error messages go to standard error. This makes the tool safe to use in scripts:

```
./tri done 3 && ./tri export | jq '.[] | select(.done)'
```

## Shell completion

The tool generates completion scripts for bash, zsh, fish, and PowerShell:

```
./tri completion zsh > ~/.zsh/completions/_tri
```

Run `./tri completion --help` for instructions for your shell.

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
