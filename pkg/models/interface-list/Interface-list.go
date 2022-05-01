package interfacelist

import (
	"fmt"

	interfaceitem "github.com/c-grimshaw/gosniff/pkg/models/interface-item"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/gopacket/pcap"
)

// Model defines the structure of the component
type Model struct {
	cursor, checked int
	interfaces      []interfaceitem.Model
}

// New returns a list of network interfaces with default parameters
func New(cursor, checked string) Model {
	interfaces, err := getInterfaces()
	if err != nil {
		panic(err)
	}

	interfaceList := make([]interfaceitem.Model, len(interfaces))
	for i, iface := range interfaces {
		interfaceList[i] = interfaceitem.New(iface, cursor, checked)
	}
	return Model{
		interfaces: interfaceList,
	}
}

// View describes the string representation of the network interface list
func (m Model) View() (view string) {
	view = "Interface:"
	for _, item := range m.interfaces {
		view = lipgloss.JoinVertical(lipgloss.Left, view, item.View())
	}
	return view
}

// getInterfaces returns all host interfaces in string format
func getInterfaces() (interfaces []pcap.Interface, err error) {
	interfaces, err = pcap.FindAllDevs()
	if err != nil {
		fmt.Println("Error: No host interfaces")
		return interfaces, err
	}
	return interfaces, nil
}
