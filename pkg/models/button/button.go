package button

import (
	"fmt"

	"github.com/c-grimshaw/gosniff/pkg/style"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the button model struct
type Model struct {
	name    string
	focused bool
}

// New returns a button model with default parameters
func New(name string) Model {
	return Model{
		name: name,
	}
}

// Update contains the button's update loop, which currently does nothing
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

// View renders the button into a text string
func (m Model) View() string {
	if m.Focused() {
		return style.Focused.Render(fmt.Sprintf("[ %s ]", m.name))
	}
	return fmt.Sprintf("[ %s ]", m.name)
}

// Focused returns the focus state of the button
func (m Model) Focused() bool {
	return m.focused
}

// SetName changes the visible display of the button
func (m *Model) SetName(name string) {
	m.name = name
}

// SetFocus changes the focus state of the button
func (m *Model) SetFocus(focus bool) {
	m.focused = focus
}
