package tui

import (
	"fmt"

	"github.com/Wondrous27/s3-tui/bucket"
	"github.com/Wondrous27/s3-tui/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	nav mode = iota
	edit
	del
)

type CreatedBucketMsg struct {
	err error
}

type DeletedBucketMsg struct {
	err        error
	bucketName string
}

type Model struct {
	mode     mode
	list     list.Model
	input    textinput.Model
	quitting bool
	isSure   bool
}

/* Implement tea.Model for Model */
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.mode == del {
		return m.DisplayConfirmation()
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
			if m.mode == del {
				switch {
				case key.Matches(msg, constants.Keymap.Quit):
					m.quitting = true
					return m, tea.Quit

				case key.Matches(msg, constants.Keymap.Enter):
					bucket := m.list.SelectedItem().(bucket.Bucket)
					m.mode = nav
					if m.isSure {
						return m, deleteBucketCommand(bucket.Name)
					} else {
						return m, nil
					}

				case key.Matches(msg, constants.Keymap.Next), key.Matches(msg, constants.Keymap.Prev):
					m.isSure = !m.isSure
				}
				return m, nil
			}
			switch {
			case key.Matches(msg, constants.Keymap.Delete):
				m.mode = del

			case key.Matches(msg, constants.Keymap.Create):
				m.input.Focus()
				cmd = textinput.Blink
				cmds = append(cmds, cmd)

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

// TODO: Come back to this
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

func (m Model) DisplayConfirmation() string {
	buttonStyle := map[bool]lipgloss.Style{
		true:  constants.ActiveButtonStyle,
		false: constants.ButtonStyle,
	}
	okb := buttonStyle[m.isSure]
	cb := buttonStyle[!m.isSure]

	okButton := okb.Render("Yes")
	cancelButton := cb.Render("No")

	activeBucket := m.list.SelectedItem().(bucket.Bucket)
	msg := fmt.Sprintf("Are you sure you want to delete %s?", activeBucket.Name)
	question := lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render(msg)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	const (
		width       = 96
		columnWidth = 30
	)

	dialog := lipgloss.Place(width, 9,
		lipgloss.Center, lipgloss.Center,
		constants.DialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars(""),
		lipgloss.WithWhitespaceForeground(constants.Subtle),
	)

	return dialog
}
