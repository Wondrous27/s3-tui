package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/Wondrous27/s3-tui/object"
	"github.com/Wondrous27/s3-tui/tui/constants"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg         struct{ error }
	UpdatedObjects []object.Object
)

type editorFinishedMsg struct {
	err  error
	file *os.File
}

var cmd tea.Cmd

// Object implements tea.Model
type Object struct {
	viewport         viewport.Model
	activeBucketName string
	error            string
	paginator        paginator.Model
	objects          []object.Object
	quitting         bool
}

// Init run any intial IO on program start
func (m Object) Init() tea.Cmd {
	return nil
}

// InitObjects initialize the objectui model for your program
func InitObjects(bucketName string) *Object {
	m := Object{activeBucketName: bucketName}
	top, right, bottom, left := constants.DocStyle.GetMargin()
	m.viewport = viewport.New(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	m.viewport.Style = lipgloss.NewStyle().Align(lipgloss.Bottom)

	// init paginator
	m.paginator = paginator.New()
	m.paginator.Type = paginator.Dots
	m.paginator.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	m.paginator.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	m.objects = m.setupObjects().(UpdatedObjects)
	m.paginator.SetTotalPages(len(m.objects))
	// set content
	m.setViewportContent()
	return &m
}

func (m *Object) setupObjects() tea.Msg {
	var err error
	var objects []object.Object
	if objects, err = constants.Or.ListObjects(m.activeBucketName); err != nil {
		return errMsg{fmt.Errorf("cannot find project: %v", err)}
	}
	return UpdatedObjects(objects)
}

func (m *Object) setViewportContent() {
	var content string
	if len(m.objects) == 0 {
		content = "There are no objects for this bucket :("
	} else {
		content = object.FormatObject(m.objects[m.paginator.Page])
	}
	str, err := glamour.Render(content, "dark")
	if err != nil {
		m.error = "could not render content with glamour"
	}
	m.viewport.SetContent(str)
}

func (m Object) helpView() string {
	// TODO: use the keymaps to populate the help string
	return constants.HelpStyle("\n ↑/↓: navigate  • esc: back • e: edit object • c: create object • d: delete entry • q: quit\n")
}

// View return the text UI to be output to the terminal
func (m Object) View() string {
	if m.quitting {
		return ""
	}

	formatted := lipgloss.JoinVertical(lipgloss.Left, "\n", m.viewport.View(), m.helpView(), constants.ErrStyle(m.error), m.paginator.View())
	return constants.DocStyle.Render(formatted)
}

func (m Object) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.viewport = viewport.New(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-6)

	case errMsg:
		m.error = msg.Error()

	case editorFinishedMsg:
		if msg.err != nil {
			return m, tea.Quit
		}
		cmds = append(cmds, m.updateObjectCmd(msg.file.Name()))

	case UpdatedObjects:
		m.objects = msg
		m.paginator.SetTotalPages(len(m.objects))
		m.setViewportContent()

	case tea.KeyMsg:
		switch {
		// TODO: find a way to override h&l to get object content on demand
		// case key.Matches(msg, constants.Keymap.Next):
		// 	fallthrough
		// case key.Matches(msg, constants.Keymap.Prev):
		// 	return m, nil
		case key.Matches(msg, constants.Keymap.Edit):
			fileContent := m.objects[m.paginator.Page].Content
			keys := strings.Split(m.objects[m.paginator.Page].Key, "/")
			fileName := keys[len(keys)-1]
			return m, openEditorCmd(fileContent, fileName)
		case key.Matches(msg, constants.Keymap.Create):
			return m, nil
			// return m, openEditorCmd()
		case key.Matches(msg, constants.Keymap.Back):
			return InitBuckets()
		case key.Matches(msg, constants.Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)
	m.setViewportContent() // refresh the content on every Update call
	return m, tea.Batch(cmds...)
}
