package gosniff

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	noStyle          = lipgloss.NewStyle()
	focusedStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	viewportStyle    = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1)
)

func (m *model) View() string {
	var view strings.Builder
	buttonInputs := lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, m.submitView(), "      ", m.clearView()),
		gap(1),
		m.helpView())
	view.WriteString(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.titleView(),
			gap(1),
			m.interfaceView(),
			gap(1),
			m.filter.View(),
			gap(1),
			buttonInputs),
	)

	block := noStyle.MaxWidth(100).Render(view.String())
	viewport := noStyle.MaxWidth(100).Render(m.viewportView())
	return lipgloss.JoinHorizontal(lipgloss.Left, block, viewport)
}

func (m *model) titleView() (view string) {
	return "//GOSNIFF//"
}

func (m *model) interfaceView() (view string) {
	view = "Interface:"
	for i, choice := range m.interfaces {
		cursor := " "
		if m.focus == i && m.focusedInterfaces() {
			cursor = ">"
		}

		row := fmt.Sprintf("%s [ ] %s", cursor, choice.Description)
		if m.selected == i {
			row = focusedStyle.Render(fmt.Sprintf("%s [x] %s", cursor, choice.Description))
		}

		view = lipgloss.JoinVertical(lipgloss.Left, view, row)
		for _, addr := range choice.Addresses {
			view = lipgloss.JoinVertical(lipgloss.Left, view, placeholderStyle.Render(fmt.Sprintf("       - [%v]", addr.IP)))
		}
	}
	return view
}

func (m *model) submitView() (view string) {
	view = "[ Start ]"
	if m.recording {
		view = "[ Stop ]"
	}

	if m.focusedSubmit() {
		view = focusedStyle.Render(view)
	}
	return view
}

func (m *model) clearView() (view string) {
	view = "[ Clear ]"
	if m.focusedClear() {
		view = focusedStyle.Render(view)
	}
	return view
}

func (m model) headerView() (view string) {
	view = viewportStyle.Render("GOSNIFF - STOPPED")
	if m.recording {
		view = viewportStyle.Render("GOSNIFF - RECORDING")
	}
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(view)))
	return lipgloss.JoinHorizontal(lipgloss.Center, view, line)
}

func (m model) footerView() (view string) {
	view = viewportStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(view)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, view)
}

func (m *model) helpView() string {
	return m.help.View(m.keys)
}

func (m *model) viewportView() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func gap(n int) string {
	return strings.Repeat("\n", n)
}
