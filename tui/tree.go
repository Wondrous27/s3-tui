package tui

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/Wondrous27/s3-tui/tree"
	"github.com/Wondrous27/s3-tui/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


type UpdatedTree *Tree

type Tree struct {
	BucketName string
	Root       *tree.Node
	quitting   bool
	cursor     int
	selected   map[string]*tree.Node
	input      textinput.Model
	mode       mode
	NewObjectKey  string
}

func (f Tree) DidSelectFile(msg tea.Msg) (bool, string) {
	return false, ""
}

func (f Tree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	// TODO: Implement custom window resizing
	case tea.WindowSizeMsg:

	case editorFinishedMsg:
		cmds = append(cmds, f.createObjectCommand(msg.file.Name(), f.NewObjectKey))

	case UpdatedTree:
		f = *msg

	case tea.KeyMsg:
		if f.input.Focused() {
			if key.Matches(msg, constants.Keymap.Back) {
				f.input.SetValue("")
				f.mode = nav
				f.input.Blur()
			}

			if key.Matches(msg, constants.Keymap.Enter) {
				s3Key := f.input.Value()
				f.NewObjectKey = s3Key
				extension := filepath.Ext(s3Key)
				f.input.SetValue("")
				f.mode = nav
				f.input.Blur()
				return f, openEditorCmd("", extension)
			}

			f.input, cmd = f.input.Update(msg)
			cmds = append(cmds, cmd)
			f.input.Update(msg)
		} else {
			switch {
			case key.Matches(msg, constants.Keymap.Quit):
				f.quitting = true
				return f, tea.Quit

			case key.Matches(msg, constants.Keymap.Create):
				f.mode = create
				f.input.Focus()
				cmd = textinput.Blink

			case key.Matches(msg, constants.Keymap.Up):
				f.cursor = (f.cursor - 1 + len(f.Root.Children)) % len(f.Root.Children)
				return f, nil

			case key.Matches(msg, constants.Keymap.Down):
				f.cursor = (f.cursor + 1) % len(f.Root.Children)
				return f, nil

			case key.Matches(msg, constants.Keymap.Enter):
				curr := f.Root.Children[f.cursor]
				if !curr.IsDir {
					key := getPath(curr)
					return InitObject(f.BucketName, key)
				}
				f.Root = curr
				f.cursor = 0
				return f, nil

			case key.Matches(msg, constants.Keymap.Back):
				return InitBuckets()

			case key.Matches(msg, constants.Keymap.Prev):
				if f.Root.Name == "" {
					return InitBuckets()
				}
				f.Root = f.Root.Parent
				f.cursor = 0
				return f, nil

			case key.Matches(msg, constants.Keymap.Next):
				if f.Root.Children[f.cursor].IsDir {
					f.Root = f.Root.Children[f.cursor]
					f.cursor = 0
				}
				return f, nil
			}
		}
	}
	return f, tea.Batch(cmds...)
}

// TODO: make this prettier
func (f Tree) View() string {
	var sb strings.Builder
	for i, child := range f.Root.Children {
		cursor := " "
		isSelected := false
		if i == f.cursor {
			cursor = "> "
			isSelected = true
		}
		sb.WriteString(cursor)
		sb.WriteString(styledFileName(child.IsDir, isSelected, child.Name))
		sb.WriteString("\n\n")
	}

	sb.WriteString(constants.HelpStyle(
		"\n ↑/↓ j/h: navigate • esc: back • c: create object • q: quit\n",
	))
	if f.input.Focused() {
		// TODO: Find new style to render this
		return constants.DocStyle.Render(sb.String() + "\n" + f.input.View())
	}
	return constants.DocStyle.Render(sb.String())
}

func styledFileName(isDir, isSelected bool, name string) string {
	if isSelected {
		return constants.SelectedStyle(name)
	}
	if isDir {
		return constants.DirStyle(name)
	}
	return constants.FileStyle(name)
}

func (f Tree) Init() tea.Cmd {
	return nil
}

func InitTree(bucketName string) *Tree {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Object Key..."
	input.CharLimit = 250
	input.Width = 50

	objects, err := constants.Or.ListObjects(bucketName)
	if err != nil {
		panic(err.Error())
	}
	root := tree.NewFileTree(objects)
	return &Tree{
		BucketName: bucketName,
		Root:       root.Root,
		cursor:     0,
		input:      input,
	}
}

func (f Tree) setupTree(bucketName string) tea.Msg {
	tree := InitTree(bucketName)
	log.Printf("%v\n", tree)
	return UpdatedTree(tree)
}

// bucket/bucket.go
func getPath(n *tree.Node) string {
	curr := n
	path := []string{curr.Name}
	for curr.Parent.Name != "" {
		curr = curr.Parent
		path = append(path, curr.Name)
	}
	var p string
	for i := len(path) - 1; i >= 0; i-- {
		p += path[i] + "/"
	}
	return p[:len(p)-1]
}
