package tui

func (m Model) View() string {
	return m.list.View()
}
