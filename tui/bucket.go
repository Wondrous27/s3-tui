package tui

import (
	"fmt"

	// "github.com/charmbracelet/bubbles/key"
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
	// var cmd tea.Cmd
	// var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case tea.KeyMsg:
		if m.input.Focused() {
			// handle the case for focused
		} else {
			switch {
			case key.Matches(msg, constants.Keymap.Quit):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, constants.Keymap.Enter):
				activeProject := m.list.SelectedItem().(bucket.Bucket)
				fmt.Println("selected", activeProject)
				// entry := InitEntry(constants.Er, activeProject.ID, constants.P)
				// return entry.Update(constants.WindowSize)
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func InitBuckets() (tea.Model, tea.Cmd) {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Bucket name..."
	input.CharLimit = 250
	input.Width = 50

	// TODO: handle error
	items, _ := newBucketsList(constants.Br)

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

func newBucketsList(br *bucket.S3Repository) ([]list.Item, error) {
	buckets, err := br.GetAllBuckets()
	if err != nil {
		return nil, fmt.Errorf("cannot get buckets: %w", err)
	}
	return bucketsToItems(buckets), err
}

// func ConvertToItems[T any](buckets []T) []list.Item {
// TODO: use generics
// projectsToItems convert []model.Project to []list.Item
func bucketsToItems(buckets []bucket.Bucket) []list.Item {
	items := make([]list.Item, len(buckets))
	for i, bucket := range buckets {
		items[i] = list.Item(bucket)
	}
	return items
}
