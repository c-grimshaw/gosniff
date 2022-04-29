package gosniff

import "github.com/charmbracelet/bubbles/key"

// KeyMap contains a list of key bindings
type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Exit  key.Binding
	Enter key.Binding
	Help  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Help}, // first column
		{k.Exit, k.Enter},      // second column
	}
}

// DefaultKeyMap describes the default key bindings
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up", "left", "shift+tab"), // actual keybindings
		key.WithHelp("↑/k", "Move up"),               // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down", "right", "tab"),
		key.WithHelp("↓/j", "Move down"),
	),
	Exit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl-c", "Exit program"),
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
