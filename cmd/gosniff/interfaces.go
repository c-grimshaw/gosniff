package gosniff

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	interfaces []pcap.Interface
	cursor     int
	focusIndex int
	selected   int
	submit     int
	keys       KeyMap
	help       help.Model
	recording  bool
	textinput  textinput.Model
}

// KeyMap contains a list of key bindings
type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Exit  key.Binding
	Next  key.Binding
	Enter key.Binding
	Help  key.Binding
}

var (
	noStyle      = lipgloss.NewStyle()
	focusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	snaplen      = int32(1600)
	promisc      = false
	timeout      = pcap.BlockForever
)

// GetInterfaces returns all host interfaces in string format
func GetInterfaces() (interfaces []pcap.Interface, err error) {
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println("Error: No host interfaces")
		return interfaces, err
	}

	for _, i := range ifaces {
		if len(i.Addresses) > 0 {
			interfaces = append(interfaces, i)
		}
	}

	return interfaces, nil
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Help},    // first column
		{k.Exit, k.Enter, k.Next}, // second column
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
	Next: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
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

func NewModel() model {
	ifaces, err := GetInterfaces()
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
	return tea.EnterAltScreen
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
			if !m.textFieldFocused() {
				m.help.ShowAll = !m.help.ShowAll
			}
		case key.Matches(msg, DefaultKeyMap.Enter):
			if m.interfaceFieldFocused() {
				if m.selected == m.cursor {
					m.selected = -1
				} else {
					m.selected = m.cursor
				}
				break
			}
			if m.textFieldFocused() {
				break
			}
			if m.submitFieldFocused() {
				if !m.recording {
					m.recording = !m.recording
					go m.start()
				}
			}
		case key.Matches(msg, DefaultKeyMap.Next):
			switch msg.String() {
			case "tab":
				m.focusIndex = mod(m.focusIndex+1, NUM_ITEMS)
			case "shift+tab":
				m.focusIndex = mod(m.focusIndex-1, NUM_ITEMS)
			}
		}
	}

	// Text Input Processing
	if m.textFieldFocused() {
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
		fmt.Println(packet)
	}
}

func (m *model) interfaceFieldFocused() bool {
	return m.focusIndex == 0
}

func (m *model) textFieldFocused() bool {
	return m.focusIndex == 1
}

func (m *model) submitFieldFocused() bool {
	return m.focusIndex == 2
}

func mod(x, m int) int {
	return (x%m + m) % m
}

func (m model) View() string {
	// The header
	s := lipgloss.NewStyle().Bold(true).Underline(true).Render("//GOSNIFF//")
	s += "\n\nInterface:\n"

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
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Description)
		for _, addr := range choice.Addresses {
			s += fmt.Sprintf("\t- [%v]\n", addr.IP)
		}
	}
	s += fmt.Sprintf("\nFilter:\n %s\n\n", m.textinput.View())
	if m.focusIndex == 2 {
		s += focusedStyle.Render("[ Start ]\n")
	} else {
		s += "[ Start ]\n"
	}

	helpView := m.help.View(m.keys)
	height := 2
	return "\n" + s + strings.Repeat("\n", height) + helpView
}
