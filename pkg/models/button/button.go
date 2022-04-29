package button

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
)

type Model struct {
	name    string
	focused bool
}

func New(name string) Model {
	return Model{
		name: name,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	if m.Focused() {
		return focusedStyle.Render(fmt.Sprintf("[ %s ]", m.name))
	}
	return fmt.Sprintf("[ %s ]", m.name)
}

func (m Model) Focused() bool {
	return m.focused
}

func (m *Model) SetName(name string) {
	m.name = name
}

func (m *Model) SetFocus(focus bool) {
	m.focused = focus
}
