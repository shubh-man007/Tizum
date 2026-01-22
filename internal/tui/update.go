package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "t":
			i := m.list.Index()
			task := m.list.Items()[i].(item)
			m.db.ToggleTask(task.ID, !task.Status)
			return refreshList(m)

		case "d":
			i := m.list.Index()
			task := m.list.Items()[i].(item)
			m.db.DeleteTask(task.ID)
			return refreshList(m)

		case "a":
			// TODO: open input mode
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func refreshList(m Model) (tea.Model, tea.Cmd) {
	tasks, _ := m.db.ReadTasks()
	var items []list.Item
	for _, t := range tasks {
		items = append(items, item(t))
	}

	m.list.SetItems(items)
	return m, nil
}
