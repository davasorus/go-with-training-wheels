# Todo Module

This module handles task management logic. It manages the core data structures and business rules for tasks.

## Components

- **todo.go**: This file defines the `Todo` struct. It contains methods to format labels. It includes sorting logic based on priority and position.
- **repository.go**: This file implements the Repository pattern. It handles all database interactions. It separates SQL queries from business logic.

## Logic Flow

1. The system uses `LoadItems` to fetch tasks.
2. `LoadItems` calls the repository layer.
3. The repository executes a SQL query.
4. The result is a list of `Todo` objects.
5. The system uses `SaveItems` to save tasks.
6. It uses transactions to ensure data integrity.

## Data Structure

The `Todo` struct includes:
- **Text**: The description of the task.
- **Priority**: An integer for importance.
- **Position**: The order in a list.
- **Done**: A boolean state for completion.
