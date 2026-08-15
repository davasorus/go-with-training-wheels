# Database Module

This module manages the connection to the PostgreSQL database. It handles configuration and initialization.

## Components

- **db.go**: This file contains the `InitDB` function. It loads environment variables. It establishes a connection to the database. It creates the database if it does not exist.

## Configuration

The system uses an `.env` file in this directory for configuration.
The following variables are required:
- **DB_HOST**: The host of the database server.
- **DB_PORT**: The port of the database server.
- **DB_USER**: The username for the database connection.
- **DB_PASSWORD**: The password for the database user.
- **DB_NAME**: The name of the database.

## Logic Flow

1. The system calls `InitDB`.
2. The function reads variables from the `.env` file.
3. It constructs a Data Source Name (DSN) string.
4. It opens a connection to the PostgreSQL server.
5. If the ping fails, it checks if the database exists.
6. It creates the database if it is missing.
7. The function returns successfully once the connection is ready.
