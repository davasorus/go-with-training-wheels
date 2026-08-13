# Project Migration Plan: JSON to PostgreSQL

This document outlines the step-by-step plan for migrating the data storage from local JSON files to a PostgreSQL database, while introducing a Repository pattern to decouple business logic from persistence details.

---

## Phase 1: Database & Persistence Layer
The goal of this phase is to establish a stable connection and create an abstraction layer (Repository) so that `todo/todo.go` does not need to know about SQL query implementation details.

- [x] **Create Database Helper**: Create `Database/db.go` to manage the connection pool using `sql.DB`.
- [ ] **Configuration**: Load configuration from `.env` and ensure the database name is set to `todos`.
- [ ] **Repository Implementation**: Implement a repository pattern (e.g., `todo/repository.go`). This layer will handle all SQL queries, mapping results to the domain models.
- [ ] **Migration Setup**: Ensure the schema is initialized using the Liquibase `changelog.yaml` file.
- [ ] **Connection Test**: Add tests to ensure connectivity and successful execution of basic queries through the repository.

---

## Phase 2: Refactor Core Logic (`todo/todo.go`)
This phase involves migrating business logic from JSON manipulation to interacting with the Repository abstraction.

- [ ] **Remove JSON dependencies**: Remove `encoding/json` and file system calls from `todo/todo.go`.
- [ ] **Integrate Repository**: Update functions like `SaveItems` and `LoadItems` to use the repository methods instead of raw SQL or file reads.
- [ ] **Refactor Model Logic**: Clean up the `Todo` struct and related logic, ensuring it focuses only on domain rules (validation, formatting).

---

## Phase 3: Update Command Layer (`cmd/`)
This phase ensures that the CLI interactions are correctly mapped to the refactored logic.

- [ ] **Update Commands**: Refactor `add.go`, `list.go`, and `done.go` to interact with the high-level `todo` functions which now utilize the Repository under the hood.
- [ ] **Integration Test**: Run full end-to-end tests to ensure a task can be added, listed, and marked as done using the PostgreSQL backend via the Repository layer.
