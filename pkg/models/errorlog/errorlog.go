package errorlog

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Model is the error log model struct
type Model struct {
	content string
}

// New returns an error log model with default params
func New() Model {
	return Model{}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("%s", m.content)
}

func (m *Model) SetContent(s string) {
	m.content = s
}
