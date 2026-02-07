package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const footerKeys = "↑/k up • ↓/j down • / filter • t toggle • d delete • a add • e edit • q quit"

func (m Model) View() string {
	banner := renderBanner(m.width)
	datetime := renderCurrentDateTime(m.width)
	listView := m.list.View()

	footer := footerStyle.Render(footerKeys)
	if m.width > 0 {
		footer = lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(footerKeys)
	}

	var b strings.Builder
	b.WriteString(banner)
	b.WriteString("\n")
	b.WriteString(datetime)
	b.WriteString("\n")
	b.WriteString(listView)
	b.WriteString("\n")
	b.WriteString(footer)

	if m.state == addState {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(" New task: "))
		b.WriteString(m.input.View())
	}

	if m.state == editState {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(" Edit: "))
		b.WriteString(m.input.View())
	}

	return b.String()
}
