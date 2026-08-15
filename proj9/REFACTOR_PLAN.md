# Refactor Plan: Consolidate Database Logic

## Objective
- Move all raw SQL queries and database interactions into the `Database` package (specifically `db.go`).
- Implement an Object-Oriented approach to decouple domain logic from persistence details.
- Ensure that no other part of the application interacts with `sql.DB` or `sql.Tx` directly.

## Scope
- `Database/db.go`: Centralize all SQL execution.
- `todo/repository.go`: Refactor to use a clean interface for database operations.
- `cmd/`: Ensure commands interact only with the domain layer.

## Steps

### Step 1: Define Database Interface and Store in `Database/db.go`
- Create a `Store` struct in the `Database` package that wraps `*sql.DB`.
- Move the SQL queries from `todo/repository.go` into methods of this `Store`.
- Example methods to move:
    - `ListTodos(query string)` (or similar, tailored to requirements)
    - `SaveTodos(items []models.Todo)`
    - `UpdateStatus(id int, done bool)`
- Replace the global `DB` variable with an instance of this Store or a well-defined interface.

### Step 2: Refactor `todo/repository.go` to use the Interface
- Define an interface (e.g., `TodoStore`) in a shared package or within the `todo` package that describes the operations needed for todos.
- Update the `Repository` struct to hold a reference to this interface instead of `*sql.DB`.
- Update implementation of `ListItems`, `SaveItems`, and `UpdateItemStatus` to call methods on the interface.

### Step 3: Cleanup and Verification
- Remove any remaining direct uses of `Database.DB` in other files.
- Update tests to ensure that the new structure still produces correct outputs for all commands.
- Verify that all SQL strings are consolidated into the `Database` package.
