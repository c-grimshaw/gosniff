package gosniff

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	snaplen = int32(1600)
	promisc = false
	timeout = pcap.BlockForever
)

// GetInterfaces returns all host interfaces in string format
func GetInterfaces() (interfaces []pcap.Interface, err error) {
	interfaces, err = pcap.FindAllDevs()
	if err != nil {
		fmt.Println("Error: No host interfaces")
		return interfaces, err
	}
	return interfaces, nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case packetMsg:
		m.content += msg.String() + "\n"
		m.viewport.SetContent(m.content)
		m.viewport.GotoBottom()
		return m, waitForPacket(m.packetChan)

	case stopMsg:
		m.stopChan <- false
		m.recording = false
		return m, waitForStop(m.stopChan)

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Up):
			m.cursorUp()
		case key.Matches(msg, DefaultKeyMap.Down):
			m.cursorDown()
		case key.Matches(msg, DefaultKeyMap.Next):
			m.handleTab(msg.String())
		case key.Matches(msg, DefaultKeyMap.Enter):
			m.handleEnter()
		case key.Matches(msg, DefaultKeyMap.Help):
			if !(m.focusedFilter()) {
				m.help.ShowAll = !m.help.ShowAll
			}
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// packetMsg is a data message sent from the packet filter
type packetMsg string

func (p packetMsg) String() string { return string(p) }

// stopMsg is a control message sent to stop the packet filter
type stopMsg struct{}

// start is used to turn on the packet filter with user-specified inputs
func (m *model) start() {
	iface := m.interfaces[m.selected].Name
	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	// TODO: Return errors on prompt
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	// TODO: Return errors on prompt
	if err := handle.SetBPFFilter(m.textinput.Value()); err != nil {
		log.Panicln(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case packet := <-source.Packets():
			m.packetChan <- packet
		case <-m.stopChan:
			return
		}
	}
}

// cursorUp moves the cursor up under the interfaces input
func (m *model) cursorUp() {
	if m.focus > 0 {
		m.focus--
	}
}

// cursorDown moves the cursor down under the interfaces input
func (m *model) cursorDown() {
	if m.focus < m.inputs {
		m.focus++
	}
}

// handleTab controls tab movement through input fields
func (m *model) handleTab(key string) {
	switch key {
	case "tab":
		m.focus = mod(m.focus+1, m.inputs)
	case "shift+tab":
		m.focus = mod(m.focus-1, m.inputs)
	}

	if m.focusedFilter() {
		// Set focused state
		m.textinput.Focus()
		m.textinput.PromptStyle = focusedStyle
		m.textinput.TextStyle = focusedStyle
	} else {
		// Remove focused state
		m.textinput.Blur()
		m.textinput.PromptStyle = noStyle
		m.textinput.TextStyle = noStyle
	}
}

// handleEnter controls enter behaviour over input fields
func (m *model) handleEnter() {
	if m.focusedInterfaces() {
		m.selected = m.focus
	}
	if m.focusedFilter() {
		return
	}
	if m.focusedSubmit() {
		// Start the packet listener goroutine
		if !m.recording {
			m.recording = !m.recording
			go m.start()
		} else {
			m.stopChan <- true
		}
		return
	}
	if m.focusedClear() {
		m.content = ""
		m.viewport.SetContent(m.content)
	}
}

// waitForPacket is a listener that sends received packets to the main model for display
// in the viewport component
func waitForPacket(packet chan gopacket.Packet) tea.Cmd {
	return func() tea.Msg {
		return packetMsg((<-packet).String())
	}
}

// waitForStop is a listener that emits a stopMsg when the recording is stopped
func waitForStop(stop chan bool) tea.Cmd {
	return func() tea.Msg {
		<-stop
		return stopMsg{}
	}
}

func (m *model) focusInterfaces() { m.focus = 0 }
func (m *model) focusedInterfaces() bool {
	return m.focus >= 0 && m.focus < len(m.interfaces)
}
func (m *model) focusFilter() {
	m.focus = len(m.interfaces)
}
func (m *model) focusedFilter() bool {
	return m.focus == len(m.interfaces)
}
func (m *model) focusSubmit() {
	m.focus = len(m.interfaces) + 1
}
func (m *model) focusedSubmit() bool {
	return m.focus == len(m.interfaces)+1
}
func (m *model) focusClear() {
	m.focus = len(m.interfaces) + 2
}
func (m *model) focusedClear() bool {
	return m.focus == len(m.interfaces)+2
}
func mod(x, m int) int {
	return (x%m + m) % m
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
