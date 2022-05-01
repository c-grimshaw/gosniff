package gosniff

import (
	"fmt"

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
		m.recording = false
		m.stopChan <- struct{}{}
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
		case key.Matches(msg, DefaultKeyMap.Enter):
			m.handleEnter()
		case key.Matches(msg, DefaultKeyMap.Help):
			if !(m.filter.Focused()) {
				m.help.ShowAll = !m.help.ShowAll
			}
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.filter, cmd = m.filter.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// packetMsg is a data message sent from the packet filter
type packetMsg string

func (p packetMsg) String() string { return string(p) }

// stopMsg is a control message sent to stop the packet filter
type stopMsg struct{}

// start is used to turn on the packet filter with user-specified inputs
func (m *model) listenForPackets() {
	iface := m.interfaces[m.selected].Name
	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		m.errorLog.SetContent("Error: Invalid interface!")
		return
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(m.filter.Value()); err != nil {
		m.errorLog.SetContent("Error: Invalid filter!")
		return
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
	m.focus = mod(m.focus-1, len(m.interfaces)+3)
	m.checkFocus()
}

// cursorDown moves the cursor down under the interfaces input
func (m *model) cursorDown() {
	m.focus = mod(m.focus+1, len(m.interfaces)+3)
	m.checkFocus()
}

// handleEnter controls enter behaviour over input fields
func (m *model) handleEnter() {
	switch m.focus {

	// cursor over filter
	case m.focusedFilterIdx():
		break

	// cursor over submit
	case m.focusedSubmitIdx():
		if !m.recording {
			m.recording = true
			m.submit.SetName("Stop")
			go m.listenForPackets()
		} else {
			m.submit.SetName("Start")
			m.stopChan <- struct{}{}
		}

	// cursor over clear
	case m.focusedClearIdx():
		m.clearViewport()

	// default case: cursor over interfaces
	default:
		m.selected = m.focus
	}
}

func (m *model) clearViewport() {
	m.content = ""
	m.viewport.SetContent(m.content)
}

func (m *model) checkFocus() {
	switch m.focus {
	case m.focusedFilterIdx():
		m.filter.SetFocus(true)
		m.submit.SetFocus(false)
		m.clear.SetFocus(false)
	case m.focusedSubmitIdx():
		m.filter.SetFocus(false)
		m.submit.SetFocus(true)
		m.clear.SetFocus(false)
	case m.focusedClearIdx():
		m.filter.SetFocus(false)
		m.submit.SetFocus(false)
		m.clear.SetFocus(true)
	default:
		m.filter.SetFocus(false)
		m.submit.SetFocus(false)
		m.clear.SetFocus(false)
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
func waitForStop(stop chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-stop
		return stopMsg{}
	}
}

func (m *model) focusedFilterIdx() int { return len(m.interfaces) }
func (m *model) focusedSubmitIdx() int { return len(m.interfaces) + 1 }
func (m *model) focusedClearIdx() int  { return len(m.interfaces) + 2 }

func (m *model) focusedInterfaces() bool { return m.focus < len(m.interfaces) }
func (m *model) focusedSubmit() bool     { return m.focus == len(m.interfaces)+1 }
func (m *model) focusedClear() bool      { return m.focus == len(m.interfaces)+2 }

func mod(x, m int) int {
	return (x%m + m) % m
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
