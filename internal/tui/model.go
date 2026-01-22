package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shubh-man007/Tizu/internal/models"
	"github.com/shubh-man007/Tizu/internal/repository"
)

type Model struct {
	list list.Model
	db   *repository.TizOrch
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

	l := list.New(items, list.NewDefaultDelegate(), 40, 10)
	l.Title = "Your Tasks"

	return Model{list: l, db: db}
}

func Run(m Model) (tea.Model, error) {
	return tea.NewProgram(m).Run()
}
