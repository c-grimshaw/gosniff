package gosniff

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()

	noStyle          = lipgloss.NewStyle()
	focusedStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func (m model) View() string {
	titleBlock := "//GOSNIFF//\n"
	s := "Interface:"
	s = lipgloss.JoinVertical(lipgloss.Left, titleBlock, s)

	s = lipgloss.JoinVertical(lipgloss.Left, s, m.interfaceView())
	s = lipgloss.JoinVertical(lipgloss.Left, s, m.filterView())
	s = lipgloss.JoinVertical(lipgloss.Center, s, m.submitView())

	helpView := m.help.View(m.keys)
	s = lipgloss.JoinVertical(lipgloss.Left, s, "\n"+helpView)
	block := lipgloss.NewStyle().MaxWidth(100).Render(s)
	block2 := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	block2 = lipgloss.NewStyle().MaxWidth(100).Render(block2)
	return lipgloss.JoinHorizontal(lipgloss.Left, block, block2)
}

func (m *model) interfaceView() (view string) {
	for i, choice := range m.interfaces {
		cursor := " "
		if m.cursor == i && m.focusIndex == 0 {
			cursor = ">"
		}

		checked := " "
		description := choice.Description
		if i == m.selected {
			checked = "x"
		}

		view = lipgloss.JoinVertical(lipgloss.Left, view, fmt.Sprintf("%s [%s] %s", cursor, checked, description))
		for _, addr := range choice.Addresses {
			view = lipgloss.JoinVertical(lipgloss.Left, view, placeholderStyle.Render(fmt.Sprintf("       - [%v]", addr.IP)))
		}
	}
	return view
}

func (m *model) filterView() (view string) {
	return fmt.Sprintf("Filter:\n %s", m.textinput.View())
}

func (m *model) submitView() (view string) {
	view = "[ Start ]"
	if m.focusIndex == submitInput {
		view = focusedStyle.Render("[ Start ]")
	}
	return view
}

func (m model) headerView() string {
	title := titleStyle.Render("GOSNIFF - STOPPED")
	if m.recording {
		title = titleStyle.Render("GOSNIFF - RECORDING")
	}
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
