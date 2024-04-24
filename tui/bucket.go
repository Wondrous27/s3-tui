package tui

import (

	// "github.com/charmbracelet/bubbles/key"
	"log"

	"github.com/Wondrous27/s3-tui/bucket"
	"github.com/Wondrous27/s3-tui/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

/* TODO */
// type (
// 	updateBucketListMsg struct{}
// 	renameProjectMsg     []list.Item
// )

type mode int

// TODO: When deleting a bucket
// question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
const (
	nav mode = iota
	edit
	create // TODO: create bucket - aws s3 mb s3://bucket-name
	rename // TODO: rename bucket - aws s3 mb s3://[new-bucket] && aws s3 sync s3://[old-bucket] s3://[new-bucket] && aws s3 rb --force s3://[old-bucket]
)

type CreatedBucketMsg struct {
	err error
}

type DeletedBucketMsg struct {
	err error
}

type Model struct {
	mode     mode
	list     list.Model
	input    textinput.Model
	quitting bool
}

/* Implement tea.Model for Model */
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.input.Focused() {
		return constants.DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}

	return constants.DocStyle.Render(m.list.View() + "\n")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)

	case CreatedBucketMsg:
		m.setupBuckets()

	case DeletedBucketMsg:
		m.setupBuckets()

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		if m.input.Focused() {
			if key.Matches(msg, constants.Keymap.Back) {
				m.input.SetValue("")
				m.mode = nav
				m.input.Blur()
			}

			if key.Matches(msg, constants.Keymap.Enter) {
				bucketName := m.input.Value()
				m.input.SetValue("")
				m.mode = nav
				m.input.Blur()
				return m, createBucketCommand(bucketName)
			}

			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
			m.input.Update(msg)
		} else {
			switch {
			case key.Matches(msg, constants.Keymap.Delete):
				bucket := m.list.SelectedItem().(bucket.Bucket)
				log.Println("deleting bucket:", bucket.Name)
				return m, deleteBucketCommand(bucket.Name)

			case key.Matches(msg, constants.Keymap.Create):
				m.mode = create
				m.input.Focus()
				cmd = textinput.Blink

			case key.Matches(msg, constants.Keymap.Quit):
				m.quitting = true
				return m, tea.Quit

			case key.Matches(msg, constants.Keymap.Enter), key.Matches(msg, constants.Keymap.Next):
				activeBucket := m.list.SelectedItem().(bucket.Bucket)
				tree := InitTree(activeBucket.Name)
				return tree.Update(constants.WindowSize)
			}
		}
	}

	tea.Batch(cmds...)
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// TODO: This is inefficient, call this function only for the first time. Use setupBuckets() for subsequent calls.
func InitBuckets() (tea.Model, tea.Cmd) {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Bucket name..."
	input.CharLimit = 250
	input.Width = 50

	items, err := constants.Br.GetAllBuckets()
	if err != nil {
		return nil, func() tea.Msg {
			return errMsg{error: err}
		}
	}

	m := Model{mode: nav, list: list.New(items, list.NewDefaultDelegate(), 8, 8), input: input}
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}

	m.list.Title = "buckets"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Create,
			constants.Keymap.Rename,
			constants.Keymap.Delete,
			constants.Keymap.Back,
		}
	}
	return m, nil
}

func (m *Model) setupBuckets() tea.Msg {
	items, err := constants.Br.GetAllBuckets()
	if err != nil {
		return errMsg{error: err}
	}
	m.list.SetItems(items)
	return nil
}
