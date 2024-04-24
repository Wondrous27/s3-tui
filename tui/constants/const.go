package constants

import (
	"github.com/Wondrous27/s3-tui/bucket"
	"github.com/Wondrous27/s3-tui/object"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// P the current tea program
	P *tea.Program
	// Br the bucket repository for the tui
	Br *bucket.S3Repository
	// Or the object repository for the tui
	Or *object.S3Repository
	// WindowSize store the size of the terminal window
	WindowSize tea.WindowSizeMsg
)

/* STYLING */

// DocStyle styling for viewports
var DocStyle = lipgloss.NewStyle().Margin(0, 2)

// HelpStyle styling for help context menu
var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

// ErrStyle provides styling for error messages
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd534b")).Render

// AlertStyle provides styling for alert messages
var AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render

var (
	DirStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render
	FileStyle = lipgloss.NewStyle().Bold(true).Render
	// Selected Style background grey color
	SelectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("241")).Render
	// FileStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("23")).Render
)

var (
	ButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginRight(2).
			MarginTop(1)

	ActiveButtonStyle = ButtonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)

var Subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

type keymap struct {
	Create key.Binding
	Edit   key.Binding
	Enter  key.Binding
	Rename key.Binding
	Delete key.Binding
	Back   key.Binding
	Quit   key.Binding
	Next   key.Binding
	Prev   key.Binding
	Up     key.Binding
	Down   key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
	),
	Next: key.NewBinding(
		key.WithKeys("l"),
	),
	Prev: key.NewBinding(
		key.WithKeys("h"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
	),
}
