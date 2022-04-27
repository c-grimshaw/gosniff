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
	view.WriteString(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.titleView(),
			gap(1),
			m.interfaceView(),
			gap(1),
			m.filterView(),
			gap(1),
			m.submitView(),
			gap(1),
			m.helpView()),
	)

	block := lipgloss.NewStyle().MaxWidth(100).Render(view.String())
	block2 := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	block2 = lipgloss.NewStyle().MaxWidth(100).Render(block2)
	return lipgloss.JoinHorizontal(lipgloss.Left, block, block2)
}

func (m *model) titleView() (view string) {
	return "//GOSNIFF//"
}

func (m *model) interfaceView() (view string) {
	view = "Interface:"
	for i, choice := range m.interfaces {
		cursor := " "
		if m.cursor == i && m.focusIndex == interfaceInput {
			cursor = ">"
		}

		row := ""
		if m.selected == i {
			row = fmt.Sprintf("%s [x] %s", cursor, choice.Description)
			row = focusedStyle.Render(row)
		} else {
			row = fmt.Sprintf("%s [ ] %s", cursor, choice.Description)
		}
		view = lipgloss.JoinVertical(lipgloss.Left, view, row)
		for _, addr := range choice.Addresses {
			view = lipgloss.JoinVertical(lipgloss.Left, view, placeholderStyle.Render(fmt.Sprintf("       - [%v]", addr.IP)))
		}
	}
	return view
}

func (m *model) filterView() (view string) {
	view = "Filter:"
	if len(m.textinput.Value()) > 0 {
		return lipgloss.JoinVertical(lipgloss.Left, view, focusedStyle.Render(m.textinput.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Left, view, m.textinput.View())
}

func (m *model) submitView() (view string) {
	view = "[ Start ]"
	if m.recording {
		view = "[ Stop ]"
	}

	if m.focusIndex == submitInput {
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

func (m model) footerView() string {
	info := viewportStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
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
