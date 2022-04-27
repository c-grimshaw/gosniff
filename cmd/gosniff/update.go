package gosniff

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const INPUTS = 3
const (
	interfaceInput = iota
	filterInput
	submitInput
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
		m.viewport, cmd = m.viewport.Update(msg)
		m.viewport.GotoBottom()
		cmds = append(cmds, cmd)

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
			if !(m.focusIndex == filterInput) {
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

type packetMsg string

func (p packetMsg) String() string { return string(p) }

func Process(packet gopacket.Packet) tea.Msg { return packetMsg(packet.String()) }

func (m *model) start() {
	iface := m.interfaces[m.selected].Name
	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(m.textinput.Value()); err != nil {
		log.Panicln(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		m.Update(Process(packet))
	}
}

func (m *model) cursorUp() {
	if m.cursor > 0 && m.focusIndex == 0 {
		m.cursor--
	}
}

func (m *model) cursorDown() {
	if m.cursor < len(m.interfaces)-1 && m.focusIndex == 0 {
		m.cursor++
	}
}

func (m *model) handleTab(key string) {
	switch key {
	case "tab":
		m.focusIndex = mod(m.focusIndex+1, INPUTS)
	case "shift+tab":
		m.focusIndex = mod(m.focusIndex-1, INPUTS)
	}

	if m.focusIndex == filterInput {
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

func (m *model) handleEnter() {
	if m.focusIndex == interfaceInput {
		if m.selected == m.cursor {
			m.selected = -1
		} else {
			m.selected = m.cursor
		}
	}
	if m.focusIndex == submitInput && !m.recording {
		m.recording = !m.recording
		go m.start()
	}
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
