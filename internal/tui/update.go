package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		bannerHeight := 7 // banner (6) + current datetime (1)
		footerHeight := 1
		inputHeight := 0
		if m.state == addState || m.state == editState {
			inputHeight = 1
		}
		listHeight := m.height - bannerHeight - footerHeight - inputHeight
		if listHeight < 4 {
			listHeight = 4
		}
		m.list.SetSize(msg.Width, listHeight)
		return m, nil

	case tea.KeyMsg:
		switch m.state {

		case addState:
			switch msg.String() {
			case "esc":
				m.state = listState
				m.input.Reset()
				m.input.Blur()
				return m, nil
			case "enter":
				text := m.input.Value()
				if text == "" {
					return m, nil
				}
				_ = m.db.CreateTask(text)
				m.input.Reset()
				m.state = listState
				m.input.Blur()
				var cmd tea.Cmd
				m, cmd = refreshList(m)
				return m, cmd
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd

		case editState:
			switch msg.String() {
			case "esc":
				m.state = listState
				m.input.Reset()
				m.input.Blur()
				m.editingID = 0
				return m, nil
			case "enter":
				text := m.input.Value()
				if text == "" || m.editingID == 0 {
					return m, nil
				}
				_ = m.db.EditTask(m.editingID, text)
				m.input.Reset()
				m.state = listState
				m.input.Blur()
				m.editingID = 0
				var cmd tea.Cmd
				m, cmd = refreshList(m)
				return m, cmd
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd

		case listState:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "t":
				if len(m.list.Items()) == 0 {
					return m, nil
				}
				i := m.list.Index()
				task := m.list.Items()[i].(item)
				_ = m.db.ToggleTask(task.ID, !task.Status)
				return refreshList(m)
			case "d":
				if len(m.list.Items()) == 0 {
					return m, nil
				}
				i := m.list.Index()
				task := m.list.Items()[i].(item)
				_ = m.db.DeleteTask(task.ID)
				return refreshList(m)
			case "a":
				m.state = addState
				m.input.Reset()
				m.input.Focus()
				return m, textinput.Blink
			case "e":
				if len(m.list.Items()) == 0 {
					return m, nil
				}
				i := m.list.Index()
				task := m.list.Items()[i].(item)
				m.state = editState
				m.editingID = task.ID
				m.input.SetValue(task.TaskName)
				m.input.Focus()
				return m, textinput.Blink
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func refreshList(m Model) (Model, tea.Cmd) {
	tasks, _ := m.db.ReadTasks()
	var items []list.Item
	for _, t := range tasks {
		items = append(items, item(t))
	}
	m.list.SetItems(items)
	return m, nil
}
