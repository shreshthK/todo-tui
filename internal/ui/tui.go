package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shreshthkandari/todo-tui/internal/core"
	"github.com/shreshthkandari/todo-tui/internal/store"
)

// StartTUI boots the terminal UI and blocks until the user quits.
func StartTUI(taskStore *store.JSONStore, initialTasks []core.Task) error {
	tui := &todoTUI{
		app:   tview.NewApplication(),
		store: taskStore,
		tasks: append([]core.Task(nil), initialTasks...),
	}
	tui.buildLayout()
	tui.refreshList(0)
	tui.setStatus("Ready")
	return tui.app.SetRoot(tui.pages, true).EnableMouse(true).Run()
}

type todoTUI struct {
	app    *tview.Application
	pages  *tview.Pages
	list   *tview.List
	status *tview.TextView

	store *store.JSONStore
	tasks []core.Task
}

func (t *todoTUI) buildLayout() {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[::b]Todo TUI[::-]  [gray]Ctrl+A add  Space toggle  Ctrl+D delete  ? help  Q quit")
	header.SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	t.list = tview.NewList().
		ShowSecondaryText(true)
	t.list.SetBorder(true).SetTitle(" Tasks ")
	t.list.SetMainTextStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite))
	t.list.SetSecondaryTextColor(tcell.ColorGray)
	t.list.SetSelectedTextColor(tcell.ColorBlack)
	t.list.SetSelectedBackgroundColor(tcell.ColorLightCyan)
	t.list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		task, ok := t.selectedTask(index)
		if !ok {
			return
		}
		state := "pending"
		if task.Done {
			state = "done"
		}
		t.setStatus(fmt.Sprintf("Selected #%d (%s)", task.ID, state))
	})

	t.status = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	t.status.SetBorder(true).SetTitle(" Status ")

	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(t.list, 0, 1, true).
		AddItem(t.status, 3, 0, false)

	t.pages = tview.NewPages().
		AddPage("main", content, true, true)

	t.app.SetInputCapture(t.handleKeys)
}

func (t *todoTUI) handleKeys(event *tcell.EventKey) *tcell.EventKey {
	front, _ := t.pages.GetFrontPage()
	onMainPage := front == "main"

	switch {
	case event.Key() == tcell.KeyCtrlC:
		t.app.Stop()
		return nil
	case onMainPage && event.Key() == tcell.KeyRune && (event.Rune() == 'q' || event.Rune() == 'Q'):
		t.app.Stop()
		return nil
	case onMainPage && event.Key() == tcell.KeyCtrlA:
		t.showAddModal()
		return nil
	case onMainPage && event.Key() == tcell.KeyCtrlD:
		t.deleteSelectedTask()
		return nil
	case onMainPage && event.Key() == tcell.KeyRune && event.Rune() == ' ':
		t.toggleSelectedDone()
		return nil
	case onMainPage && event.Key() == tcell.KeyRune && event.Rune() == '?':
		t.showHelpModal()
		return nil
	}

	return event
}

func (t *todoTUI) showAddModal() {
	input := tview.NewInputField().
		SetLabel("Task title: ").
		SetFieldWidth(40)

	form := tview.NewForm().
		AddFormItem(input).
		AddButton("Add", func() {
			title := strings.TrimSpace(input.GetText())
			updated, err := core.AddTask(t.tasks, title)
			if err != nil {
				t.setStatus(err.Error())
				return
			}
			t.tasks = updated
			if err := t.store.Save(t.tasks); err != nil {
				t.setStatus(fmt.Sprintf("save failed: %v", err))
				return
			}
			t.removeModal()
			t.refreshList(t.tasks[len(t.tasks)-1].ID)
			t.setStatus("Task added")
		}).
		AddButton("Cancel", func() {
			t.removeModal()
			t.setStatus("Add cancelled")
		})
	form.SetBorder(true).SetTitle(" Add Task ").SetTitleAlign(tview.AlignLeft)

	t.showModal("add", centered(60, 9, form))
}

func (t *todoTUI) showHelpModal() {
	help := tview.NewModal().
		SetText("Shortcuts\n\nCtrl+A: Add task\nSpace: Toggle selected task\nCtrl+D: Delete selected task\nQ: Quit").
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			t.removeModal()
		})
	help.SetTitle(" Help ").SetBorder(true)
	t.showModal("help", help)
}

func (t *todoTUI) deleteSelectedTask() {
	task, ok := t.selectedTask(t.list.GetCurrentItem())
	if !ok {
		t.setStatus("No task selected")
		return
	}

	updated, err := core.DeleteTask(t.tasks, task.ID)
	if err != nil {
		t.setStatus(err.Error())
		return
	}
	t.tasks = updated
	if err := t.store.Save(t.tasks); err != nil {
		t.setStatus(fmt.Sprintf("save failed: %v", err))
		return
	}

	nextID := task.ID
	if len(t.tasks) == 0 {
		nextID = 0
	}
	t.refreshList(nextID)
	t.setStatus(fmt.Sprintf("Deleted #%d", task.ID))
}

func (t *todoTUI) toggleSelectedDone() {
	task, ok := t.selectedTask(t.list.GetCurrentItem())
	if !ok {
		t.setStatus("No task selected")
		return
	}

	updated, err := core.ToggleDone(t.tasks, task.ID)
	if err != nil {
		t.setStatus(err.Error())
		return
	}
	t.tasks = updated
	if err := t.store.Save(t.tasks); err != nil {
		t.setStatus(fmt.Sprintf("save failed: %v", err))
		return
	}

	t.refreshList(task.ID)
	state := "pending"
	if !task.Done {
		state = "done"
	}
	t.setStatus(fmt.Sprintf("Task #%d -> %s", task.ID, state))
}

func (t *todoTUI) refreshList(selectedID int) {
	t.list.Clear()
	indexToSelect := -1
	for i, task := range t.tasks {
		mainText := fmt.Sprintf("[ ]  #%d %s", task.ID, task.Title)
		if task.Done {
			mainText = fmt.Sprintf("[gray][x][-]  [gray]#%d %s[-]", task.ID, task.Title)
		}
		secondaryText := fmt.Sprintf("Created %s", task.CreatedAt.Local().Format(time.DateTime))
		t.list.AddItem(mainText, secondaryText, 0, nil)
		if task.ID == selectedID {
			indexToSelect = i
		}
	}

	if len(t.tasks) == 0 {
		t.list.AddItem("[gray]No tasks yet. Press A to add one.[-]", "", 0, nil)
		t.list.SetCurrentItem(0)
		return
	}

	if indexToSelect < 0 {
		indexToSelect = 0
	}
	t.list.SetCurrentItem(indexToSelect)
}

func (t *todoTUI) selectedTask(index int) (core.Task, bool) {
	if index < 0 || index >= len(t.tasks) {
		return core.Task{}, false
	}
	return t.tasks[index], true
}

func (t *todoTUI) setStatus(message string) {
	t.status.SetText(fmt.Sprintf("[gray]%s[-]", message))
}

func (t *todoTUI) showModal(name string, primitive tview.Primitive) {
	t.pages.AddAndSwitchToPage(name, primitive, true)
	t.app.SetFocus(primitive)
}

func (t *todoTUI) removeModal() {
	front, _ := t.pages.GetFrontPage()
	if front != "main" {
		t.pages.RemovePage(front)
	}
	t.pages.SwitchToPage("main")
	t.app.SetFocus(t.list)
}

func centered(width int, height int, primitive tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(primitive, height, 1, true).
				AddItem(nil, 0, 1, false),
			width, 1, true,
		).
		AddItem(nil, 0, 1, false)
}
