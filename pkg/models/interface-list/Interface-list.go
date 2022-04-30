package netinterfaces

import (
	"fmt"

	"github.com/c-grimshaw/gosniff/pkg/style"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/gopacket/pcap"
)

// Model defines the structure of the component
type Model struct {
	cursor, selected int
	interfaces       []pcap.Interface
}

// New returns a list of network interfaces with default parameters
func New() Model {
	interfaces, err := GetInterfaces()
	if err != nil {
		panic(err)
	}
	return Model{
		interfaces: interfaces,
	}
}

// View describes the string representation of the network interface list
func (m Model) View() (view string) {
	view = "Interface:"
	for i, choice := range m.interfaces {
		cursor := " "
		if i == m.selected {
			cursor = ">"
		}

		row := fmt.Sprintf("%s [ ] %s", cursor, choice.Description)
		if m.selected == i {
			row = style.Focused.Render(fmt.Sprintf("%s [x] %s", cursor, choice.Description))
		}

		view = lipgloss.JoinVertical(lipgloss.Left, view, row)
		for _, addr := range choice.Addresses {
			view = lipgloss.JoinVertical(lipgloss.Left, view, style.Placeholder.Render(fmt.Sprintf("       - [%v]", addr.IP)))
		}
	}
	return view
}

// GetInterfaces returns all host interfaces in string format
func GetInterfaces() (interfaces []pcap.Interface, err error) {
	interfaces, err = pcap.FindAllDevs()
	if err != nil {
		fmt.Println("Error: No host interfaces")
		return interfaces, err
	}
	return interfaces, nil
}
