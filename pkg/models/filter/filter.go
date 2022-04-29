package filter

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
)

type Model struct {
	textinput textinput.Model
}

func New() Model {
	textinput := textinput.New()
	textinput.Placeholder = "tcp and port 80"
	textinput.CharLimit = 156
	textinput.Width = 40

	return Model{textinput: textinput}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	view := "Filter:"
	if m.Focused() || len(m.textinput.Value()) > 0 {
		return lipgloss.JoinVertical(lipgloss.Left, view, focusedStyle.Render(m.textinput.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Left, view, m.textinput.View())
}

func (m Model) Focused() bool {
	return m.textinput.Focused()
}
