# Todo TUI

## Project Goal
Build a simple to-do list application as a terminal UI (TUI) written in Go.

## Requirements
- Allow users to add and delete tasks
- Allow users to mark tasks as completed

## CLI Usage (Examples)
```bash
# Add a task
todo add "task text"

# List tasks
todo list

# Mark a task as done
todo done <id>

# Delete a task
todo delete <id>

# Optional filters
todo list --done
todo list --pending
```

## Data Storage
- Store tasks locally in a JSON file on the filesystem

## Tech Stack
- Language: Go
- Interface: TUI
- Base command: `todo`

## TUI Library Options
- `bubbletea` — modern, popular, and well-documented; great for richer TUI UX.
- `tview` — stable, widget-based, and productive for classic terminal layouts.
- `gocui` — lightweight and minimal if you want to keep dependencies small.
- `termui` — quick to assemble dashboards and simple widgets.

## Simple Next Ideas
- Task priority (Low/Med/High)
- Basic filters for done vs pending
- Search by text
- Archive completed tasks
