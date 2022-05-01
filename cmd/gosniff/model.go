package gosniff

import (
	"fmt"
	"os"

	"github.com/c-grimshaw/gosniff/pkg/models/button"
	"github.com/c-grimshaw/gosniff/pkg/models/errorlog"
	"github.com/c-grimshaw/gosniff/pkg/models/filter"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type model struct {
	interfaces              []pcap.Interface
	selected, focus, inputs int
	recording               bool
	content                 string
	keys                    KeyMap
	help                    help.Model
	submit                  button.Model
	clear                   button.Model
	filter                  filter.Model
	errorLog                errorlog.Model
	viewport                viewport.Model
	packetChan              chan (gopacket.Packet)
	stopChan                chan (struct{})
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

	submit, clear := button.New("Start"), button.New("Clear")

	help := help.New()
	help.ShowAll = true

	return &model{
		interfaces: interfaces,
		keys:       DefaultKeyMap,
		help:       help,
		submit:     submit,
		clear:      clear,
		errorLog:   errorlog.New(),
		filter:     filter.New(),
		viewport:   viewport.New(80, 30),
		packetChan: make(chan gopacket.Packet),
		stopChan:   make(chan struct{}),
	}
}
