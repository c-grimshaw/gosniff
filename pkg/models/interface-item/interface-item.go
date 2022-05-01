package interfaceitem

import (
	"fmt"

	"github.com/c-grimshaw/gosniff/pkg/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/gopacket/pcap"
)

// Model is the error log model struct
type Model struct {
	id                    int
	name, cursor, checked string
	addresses             []string
	focused               bool
}

// New returns an error log model with default params
func New(i pcap.Interface, cursor, checked string) Model {
	return Model{
		name:      i.Name,
		addresses: getAddresses(i.Addresses),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var view string
	if m.focused {
		view = style.Focused.Render(fmt.Sprintf("%s [%s] %s", m.cursor, m.checked, m.name))
	} else {
		view = fmt.Sprintf("%s [%s] %s", m.cursor, " ", m.name)
	}

	for _, addr := range m.addresses {
		view = lipgloss.JoinVertical(lipgloss.Left,
			view,
			style.Placeholder.Render(fmt.Sprintf("       - [%s]", addr)),
		)
	}
	return view
}

func (m *Model) SetFocus(state bool) {
	m.focused = state
}

// getAddresses returns all IPs associated with a given interface
func getAddresses(addrs []pcap.InterfaceAddress) []string {
	addresses := make([]string, len(addrs))
	for i, addr := range addrs {
		addresses[i] = addr.IP.String()
	}
	return addresses
}
