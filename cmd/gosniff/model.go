package gosniff

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// Including filter, submit, and clear
const NUM_INPUTS = 3

type model struct {
	interfaces              []pcap.Interface
	selected, focus, inputs int
	recording               bool
	content                 string
	keys                    KeyMap
	help                    help.Model
	textinput               textinput.Model
	viewport                viewport.Model
	packetChan              chan (gopacket.Packet)
	stopChan                chan (bool)
}

// Init contains initial I/O commands executed by the model
func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForPacket(m.packetChan),
		waitForStop(m.stopChan),
	)
}

// NewModel returns a gosniff model with default parameters
func NewModel() *model {
	interfaces, err := GetInterfaces()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	ti := textinput.New()
	ti.Placeholder = "tcp and port 80"
	ti.CharLimit = 156
	ti.Width = 40

	help := help.New()
	help.ShowAll = true

	return &model{
		interfaces: interfaces,
		keys:       DefaultKeyMap,
		help:       help,
		inputs:     len(interfaces) + NUM_INPUTS,
		textinput:  ti,
		viewport:   viewport.New(80, 30),
		packetChan: make(chan gopacket.Packet),
		stopChan:   make(chan bool),
	}
}
