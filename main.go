package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/c-grimshaw/gosniff/cmd/gosniff"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const NUM_ITEMS = 3

type model struct {
	interfaces []string
	cursor     int
	focusIndex int
	selected   int
	submit     int
	keys       KeyMap
	help       help.Model
	textinput  textinput.Model
}

// KeyMap contains a list of key bindings
type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Exit  key.Binding
	Next  key.Binding
	Prev  key.Binding
	Enter key.Binding
	Help  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Exit, k.Enter, k.Next}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Help}, // first column
		{k.Exit, k.Enter},      // second column
		{k.Next, k.Prev},       // third column
	}
}

// DefaultKeyMap is the default key bindings
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),        // actual keybindings
		key.WithHelp("↑/k", "Move up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "Move down"),
	),
	Exit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl-c", "Exit program"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift-tab", "Previous field"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Next field"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle help"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/spacebar", "Check/Uncheck box"),
	),
}

func newModel() model {
	ifaces, err := gosniff.GetInterfaces()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	ti := textinput.New()
	ti.Placeholder = "tcp and port 80"
	ti.CharLimit = 156
	ti.Width = 20

	help := help.New()
	help.ShowAll = true

	return model{
		interfaces: ifaces,
		keys:       DefaultKeyMap,
		help:       help,
		selected:   -1,
		submit:     2,
		textinput:  ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Down):
			if m.cursor < len(m.interfaces)-1 && m.focusIndex == 0 {
				m.cursor++
			}
		case key.Matches(msg, DefaultKeyMap.Up):
			if m.cursor > 0 && m.focusIndex == 0 {
				m.cursor--
			}
		case key.Matches(msg, DefaultKeyMap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, DefaultKeyMap.Enter):
			if m.textFieldIsFocused() {
				break
			}
			if m.focusIndex == m.submit {
				m.start()
			}
			m.selected = m.cursor
		case key.Matches(msg, DefaultKeyMap.Next):
			m.focusIndex = mod(m.focusIndex+1, NUM_ITEMS)
		case key.Matches(msg, DefaultKeyMap.Prev):
			m.focusIndex = mod(m.focusIndex-1, NUM_ITEMS)
		}
	}

	// Text Input Processing
	if m.textFieldIsFocused() {
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
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func mod(x, m int) int {
	return (x%m + m) % m
}

func (m *model) start() {
	iface = m.interfaces[m.selected]
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
		fmt.Println(packet)
	}
}

func (m *model) textFieldIsFocused() bool {
	return m.focusIndex == 1
}

func (m model) View() string {
	// The header
	s := "Select Interface:\n"

	// Iterate over our choices
	for i, choice := range m.interfaces {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i && m.focusIndex == 0 {
			cursor = focusedStyle.Render(">") // cursor
		}

		// Is this choice selected?
		checked := " " // not selected
		if i == m.selected {
			checked = focusedStyle.Render("x")
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += fmt.Sprintf("\nFilter:\n %s\n\n", m.textinput.View())
	if m.focusIndex == 2 {
		s += focusedStyle.Render("[ Submit ]\n")
	} else {
		s += "[ Submit ]\n"
	}

	helpView := m.help.View(m.keys)
	height := 2
	return "\n" + s + strings.Repeat("\n", height) + helpView
}

var (
	noStyle      = lipgloss.NewStyle()
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	iface        = "lo0"
	snaplen      = int32(1600)
	promisc      = false
	timeout      = pcap.BlockForever
	filter       = "tcp and port 80"
	devFound     = false
)

func main() {

	p := tea.NewProgram(newModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
