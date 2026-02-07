package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shubh-man007/Tizu/internal/models"
	"github.com/shubh-man007/Tizu/internal/repository"
)

type state int

const (
	listState state = iota
	addState
	editState
)

type Model struct {
	list     list.Model
	db       *repository.TizOrch
	width    int
	height   int
	state    state
	input    textinput.Model
	editingID int // task ID when in editState
}

type item models.Task

func (i item) Title() string       { return i.TaskName }
func (i item) Description() string { return i.CreatedAt.Format("Jan 2 15:04") }
func (i item) FilterValue() string { return i.TaskName }

func NewModel(db *repository.TizOrch) Model {
	tasks, _ := db.ReadTasks()
	items := make([]list.Item, len(tasks))
	for i, v := range tasks {
		items[i] = item(v)
	}

	delegate := taskDelegate{}
	l := list.New(items, delegate, 80, 20)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowFilter(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	ti := textinput.New()
	ti.Placeholder = "Task description..."
	ti.CharLimit = 500
	ti.Width = 60

	return Model{
		list:   l,
		db:     db,
		width:  80,
		height: 24,
		state:  listState,
		input:  ti,
	}
}

func Run(m Model) (tea.Model, error) {
	return tea.NewProgram(m, tea.WithAltScreen()).Run()
}
