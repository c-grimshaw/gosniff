package filter

import (
	"github.com/c-grimshaw/gosniff/pkg/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the filter model struct
type Model struct {
	textinput textinput.Model
	focused   bool
}

// New returns a filter model with default parameters
func New() Model {
	ti := textinput.New()
	ti.Placeholder = "tcp and port 80"
	ti.CharLimit = 156
	ti.Width = 40

	return Model{textinput: ti}
}

// Init contains commands that are executed upon model initialization
func (m Model) Init() tea.Cmd {
	return m.textinput.SetCursorMode(textinput.CursorBlink)
}

// Update contains the filter's update loop, which currently checks for focus
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.Focused() {
		m.textinput.Focus()
		m.textinput.PromptStyle = style.Focused
		m.textinput.TextStyle = style.Focused
	} else {
		m.textinput.Blur()
		m.textinput.PromptStyle = style.None
		m.textinput.TextStyle = style.None
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

// View renders the filter into a text string
func (m Model) View() string {
	view := "Filter:"
	if m.Focused() || len(m.Value()) > 0 {
		return lipgloss.JoinVertical(lipgloss.Left, view, style.Focused.Render(m.textinput.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Left, view, m.textinput.View())
}

// Value returns the content of the filter as a string
func (m Model) Value() string {
	return m.textinput.Value()
}

// SetFocus sets the focus state of the model
func (m *Model) SetFocus(state bool) {
	m.focused = state
}

// Focused returns the focus state of the model
func (m Model) Focused() bool {
	return m.focused
}
