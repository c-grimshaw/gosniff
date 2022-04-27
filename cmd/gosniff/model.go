package gosniff

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/gopacket/pcap"
)

type model struct {
	interfaces []pcap.Interface
	cursor     int
	focusIndex int
	selected   int
	submit     int
	recording  bool
	content    string
	keys       KeyMap
	help       help.Model
	textinput  textinput.Model
	viewport   viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
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
		selected:   -1,
		submit:     submitInput,
		textinput:  ti,
		viewport:   viewport.New(80, 30),
	}
}
