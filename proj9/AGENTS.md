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

You have a knowledge vault connected (the `vault` server / `vault_*` tools).
It is your long-term memory across sessions. Use it deliberately:

## Before starting significant work
- Search the vault first (`vault_search` or the vault server's search) for
  prior decisions, gotchas, or context on this topic. Past-you may have
  already solved this or recorded why an approach failed.
- If you find a relevant note, read it before proceeding. Do not re-derive
  what is already recorded.

## While working
- Treat the vault as authoritative for decisions and conventions the user
  has recorded. If the vault and your assumptions conflict, the vault wins —
  or ask.

## When finishing significant work
Record durable conclusions as a vault note — not routine steps, but the
things future-you would want: a decision and its rationale, a non-obvious
gotcha and its fix, an architecture choice, a dead end and why it failed.

Write notes with this frontmatter so they are queryable:

```
---
type: decision        # or: gotcha, architecture, rca, reference
project: <project>    # e.g. proj9, ess-patching
date: <YYYY-MM-DD>
tags: [<topic>, ...]
---
```

Then a short title, the conclusion first, then the reasoning. One idea per
note. Prefer patching an existing note (`vault_patch` under a heading) over
creating a near-duplicate.

## What NOT to record
- Routine tool output, file contents, or step-by-step logs (the audit trail
  captures activity automatically).
- Anything transient or already obvious from the code.

The test: would this note save a future session real time or prevent a
repeated mistake? If not, don't write it.