# Implementation Plan (Learning-Friendly)

This plan breaks the project into small, clear steps so you can learn Go while building the app.

## 0) Project Setup
- Create a new Go module.
- Decide your app name and module path (example: `todo-tui`).
- Confirm the base command will be `todo`.

## 1) Define The Data Model
- Decide what a task looks like.
- Suggested fields:
  - `ID` (int)
  - `Title` (string)
  - `Done` (bool)
  - `CreatedAt` (time.Time, optional for learning)
- Write a small struct in Go to represent a task.

## 2) Choose A Storage Format
- Use a local JSON file (simple and beginner-friendly).
- Decide file name and location (example: `~/.todo/tasks.json` or `./tasks.json`).
- Write functions to:
  - Load tasks from JSON
  - Save tasks to JSON

## 3) Build The Core Task Logic
- Implement these functions with unit tests (optional but good practice):
  - AddTask(title string)
  - DeleteTask(id int)
  - MarkDone(id int)
  - ListTasks(filter string)
- Keep all logic in a `core` package.

## 4) Add A Simple CLI Interface
- Start with basic CLI commands (no TUI yet):
  - `todo add "task text"`
  - `todo list`
  - `todo done <id>`
  - `todo delete <id>`
- Use the Go standard library (`flag` or `os.Args`) to parse commands.
- Print readable output to the terminal.

## 5) TUI With `tview`
- Add `tview` (and its dependency `tcell`) to the module.
- Build a tiny "hello TUI" with `tview.NewApplication()` and a basic `TextView`.
- Confirm it runs before integrating your app.

## 6) Build The TUI Screens
- Start with a simple list view using `tview.List`:
  - Show tasks
  - Highlight the selected task
  - Indicate done vs pending (prefix like `[x]` / `[ ]`)
- Add key bindings:
  - `a` = add task
  - `d` = delete task
  - `space` = mark done
  - `q` = quit
- Add an input form for adding tasks (use `tview.Form` or `tview.InputField`).

## 7) Connect TUI To Core Logic
- When the user performs an action in the TUI:
  - Update the task list in memory
  - Save the file
  - Refresh the list view

## 8) Polish And Test
- Try a few manual scenarios:
  - Add 3 tasks, mark 1 done, delete 1
  - Restart the app and confirm tasks persist
- Clean up errors and edge cases (invalid IDs, empty input).

## 9) Optional Simple Enhancements
- Add filters (done / pending)
- Add a search bar
- Add task priority (Low/Med/High)

## Suggested File Layout
- `cmd/todo/main.go` — entry point
- `internal/core/tasks.go` — core logic
- `internal/store/json_store.go` — JSON load/save
- `internal/ui/tui.go` — TUI layer (tview)

## Learning Tips
- Keep each step small and working before moving on.
- Test the core logic before adding the TUI.
- Don’t over-optimize early; focus on getting it working.
