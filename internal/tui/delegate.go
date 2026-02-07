package tui

import (
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

const (
	checkboxDone    = "[x]"
	checkboxPending = "[ ]"
	createdDateFmt  = "Jan 2 15:04"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("63")).Bold(true)
	doneItemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedDoneStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("63")).Bold(true).Strikethrough(true)
	createdDateStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type taskDelegate struct{}

func (d taskDelegate) Height() int                               { return 1 }
func (d taskDelegate) Spacing() int                               { return 0 }
func (d taskDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d taskDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	checkbox := checkboxPending
	if i.Status {
		checkbox = checkboxDone
	}
	createdStr := i.CreatedAt.Format(createdDateFmt)
	dateRendered := createdDateStyle.Render(createdStr)
	dateWidth := lipgloss.Width(dateRendered)
	lineWidth := m.Width()
	if lineWidth < dateWidth+2 {
		lineWidth = dateWidth + 2
	}
	leftWidth := lineWidth - dateWidth - 1
	leftContent := checkbox + " " + i.TaskName
	if runewidth.StringWidth(leftContent) > leftWidth {
		leftContent = runewidth.Truncate(leftContent, leftWidth, "…")
	}
	var line string
	if index == m.Index() {
		if i.Status {
			line = selectedDoneStyle.Render(leftContent)
		} else {
			line = selectedItemStyle.Render(leftContent)
		}
	} else {
		if i.Status {
			line = doneItemStyle.Render(leftContent)
		} else {
			line = itemStyle.Render(leftContent)
		}
	}
	pad := lineWidth - lipgloss.Width(line) - dateWidth
	if pad < 0 {
		pad = 0
	}
	io.WriteString(w, line+strings.Repeat(" ", pad)+dateRendered)
}
