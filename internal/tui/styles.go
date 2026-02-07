package tui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

const bannerText = `
████████╗██╗███████╗██╗   ██╗███╗   ███╗
╚══██╔══╝██║╚══███╔╝██║   ██║████╗ ████║
   ██║   ██║  ███╔╝ ██║   ██║██╔████╔██║
   ██║   ██║ ███╔╝  ██║   ██║██║╚██╔╝██║
   ██║   ██║███████╗╚██████╔╝██║ ╚═╝ ██║
   ╚═╝   ╚═╝╚══════╝ ╚═════╝ ╚═╝     ╚═╝`

var (
	bannerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			Bold(true).
			Align(lipgloss.Center)

	datetimeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Right)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)
)

func renderBanner(width int) string {
	return bannerStyle.Width(width).Render(bannerText)
}

func renderCurrentDateTime(width int) string {
	s := time.Now().Format("Mon, Jan 2 2006 · 15:04")
	return datetimeStyle.Width(width).Render(s)
}
