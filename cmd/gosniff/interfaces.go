package gosniff

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	recording  bool
	content    string
	keys       KeyMap
	help       help.Model
	textinput  textinput.Model
	viewport   viewport.Model
}

// KeyMap contains a list of key bindings
type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Exit    key.Binding
	Next    key.Binding
	Enter   key.Binding
	Help    key.Binding
	Display key.Binding
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

var (
	noStyle          = lipgloss.NewStyle()
	focusedStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

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
	Display: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "display"),
	),
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
	ti.Width = 20

	help := help.New()
	help.ShowAll = true

	return &model{
		interfaces: interfaces,
		keys:       DefaultKeyMap,
		help:       help,
		selected:   -1,
		submit:     2,
		textinput:  ti,
		viewport:   viewport.New(80, 30),
	}
}

func (m model) Init() tea.Cmd {
	return nil
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
		case key.Matches(msg, DefaultKeyMap.Display):
			m.content += "hello world\n"
			m.viewport.SetContent(m.content)
			m.viewport, cmd = m.viewport.Update(msg)
			m.viewport.GotoBottom()
			cmds = append(cmds, cmd)
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
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)
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

func (m model) headerView() string {
	title := titleStyle.Render("GOSNIFF")
	if m.recording {
		title = titleStyle.Render("GOSNIFF - RECORDING")
	}
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m model) View() string {
	titleBlock := "//GOSNIFF//\n"
	s := "Interface:"
	s = lipgloss.JoinVertical(lipgloss.Left, titleBlock, s)

	s = lipgloss.JoinVertical(lipgloss.Left, s, m.InterfaceView(m.interfaces))

	s = lipgloss.JoinVertical(lipgloss.Left, s, fmt.Sprintf("\nFilter:\n %s\n\n", m.textinput.View()))

	if m.focusIndex == 2 {
		s = lipgloss.JoinVertical(lipgloss.Center, s, focusedStyle.Render("[ Start ]\n"))
	} else {
		s = lipgloss.JoinVertical(lipgloss.Center, s, "[ Start ]\n")
	}

	helpView := m.help.View(m.keys)
	s = lipgloss.JoinVertical(lipgloss.Left, s, "\n"+helpView)
	block := lipgloss.NewStyle().MaxWidth(100).Render(s)
	block2 := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	block2 = lipgloss.NewStyle().MaxWidth(100).Render(block2)
	return lipgloss.JoinHorizontal(lipgloss.Left, block, block2)
}

func (m *model) InterfaceView(interfaces []pcap.Interface) (view string) {
	for i, choice := range m.interfaces {
		cursor := " "
		if m.cursor == i && m.focusIndex == 0 {
			cursor = ">"
		}

		checked := " "
		description := choice.Description
		if i == m.selected {
			checked = "x"
		}

		view = lipgloss.JoinVertical(lipgloss.Left, view, fmt.Sprintf("%s [%s] %s", cursor, checked, description))
		for _, addr := range choice.Addresses {
			view = lipgloss.JoinVertical(lipgloss.Left, view, placeholderStyle.Render(fmt.Sprintf("       - [%v]", addr.IP)))
		}
	}
	return view
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
		m.focusIndex = mod(m.focusIndex+1, NUM_ITEMS)
	case "shift+tab":
		m.focusIndex = mod(m.focusIndex-1, NUM_ITEMS)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
