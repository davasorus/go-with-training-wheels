# Project Overview
This repository contains the tri tool. The tool is a command line interface. It manages records in a PostgreSQL database.

# Repository Structure
- cmd/: This directory contains CLI commands and Cobra definitions.
- Database/: This directory handles database connections and migrations.
- todo/: This directory contains task data files.
- main.go: This file is the entry point for the application.

# Build and Test
Run `go build -o tri .` to build the project.
Run `./tri` to execute the tool.
Run `go test ./...` to run all tests.

# Code Conventions
Follow standard Go style guidelines. Use Cobra patterns for CLI commands. Write short functions with one purpose. Use consistent names for variables and types.

# Commit Messages
All commits must follow the Conventional Commits standard. Use the format type(scope): description.
Allowed types are: feat, fix, docs, style, refactor, perf, test, build, ci, chore.
Use the imperative mood for descriptions.
Keep subject lines under 72 characters.
Provide a body for non-trivial changes.
Add an exclamation mark after type/scope to signal breaking changes.

# Documentation Language
Write all documentation in ASD-STE100 Simplified Technical English. Use active voice. Use simple verb tenses. Write one instruction per sentence. Keep sentences short. One term must mean one thing. Do not use vague words. Do not use -ing forms.
