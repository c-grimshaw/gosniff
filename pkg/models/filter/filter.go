package filter

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
)

type Model struct {
	textinput textinput.Model
	Focused   bool
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "tcp and port 80"
	ti.CharLimit = 156
	ti.Width = 40

	return Model{textinput: ti}
}

// Blinking isn't working?
func (m Model) Init() tea.Cmd {
	return m.textinput.SetCursorMode(textinput.CursorBlink)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.Focused {
		m.textinput.Focus()
		m.textinput.PromptStyle = focusedStyle
		m.textinput.TextStyle = focusedStyle
	} else {
		m.textinput.Blur()
		m.textinput.PromptStyle = noStyle
		m.textinput.TextStyle = noStyle
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	view := "Filter:"
	if m.Focused || len(m.Value()) > 0 {
		return lipgloss.JoinVertical(lipgloss.Left, view, focusedStyle.Render(m.textinput.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Left, view, m.textinput.View())
}

func (m Model) Value() string {
	return m.textinput.Value()
}
