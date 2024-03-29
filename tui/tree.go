package tui

import (
	"fmt"
	"strings"

	"github.com/Wondrous27/s3-tui/tree"
	"github.com/Wondrous27/s3-tui/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Tree struct {
	BucketName string
	Root       *tree.Node
	quitting   bool
	cursor     int
	selected   map[string]*tree.Node
}

func (f Tree) DidSelectFile(msg tea.Msg) (bool, string) {
	return false, ""
}

func (f Tree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// todo:
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			f.quitting = true
			return f, tea.Quit

		case key.Matches(msg, constants.Keymap.Up):
			f.cursor = (f.cursor - 1 + len(f.Root.Children)) % len(f.Root.Children)
			return f, nil

		case key.Matches(msg, constants.Keymap.Enter):
			curr := f.Root.Children[f.cursor]
			if !curr.IsDir {
				key := getPath(curr)
				fmt.Println("key: ", key, "bucket: ", f.BucketName)
				return InitObject(f.BucketName, key)
			}
			f.Root = curr
			f.cursor = 0
			return f, nil

		case key.Matches(msg, constants.Keymap.Down):
			f.cursor = (f.cursor + 1) % len(f.Root.Children)
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
	return f, nil
}

func (f Tree) View() string {
	var s strings.Builder
	for i, child := range f.Root.Children {
		cursor := " "
		isSelected := false
		if i == f.cursor {
			cursor = "> "
			isSelected = true
		}
		s.WriteString(cursor)
		s.WriteString(styledFileName(child.IsDir, isSelected, child.Name))
		s.WriteString("\n\n")
	}

	s.WriteString(constants.HelpStyle(
		"\n ↑/↓ j/h: navigate • esc: back • e: edit object • q: quit\n",
	))
	return s.String()
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
	objects, err := constants.Or.ListObjects(bucketName)
	if err != nil {
		panic(err.Error())
	}
	root := tree.NewFileTree(objects)
	return &Tree{
		BucketName: bucketName,
		Root:       root.Root,
		cursor:     0,
	}
}

// bucket/bucket.go
func getPath(n *tree.Node) string {
	curr := n
	path := []string{curr.Name}
	for curr.Parent.Name != "" {
		curr = curr.Parent
		// fmt.Println("curr.Name: ", curr.Name)
		path = append(path, curr.Name)
	}
	var p string
	for i := len(path) - 1; i >= 0; i-- {
		p += path[i] + "/"
	}
	fmt.Println("path: ", p)
	return p[:len(p)-1]
	// fmt.Println("path: ", path)
	// fmt.Println("path joined: ", strings.Join(path, "/"))
	// return strings.Join(path, "/")
}
