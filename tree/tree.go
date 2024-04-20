package tree

import (
	"sort"
	"strings"
	"time"
)

// TODO: Handle duplicate names
type Node struct {
	Name         string
	IsDir        bool
	Parent       *Node
	Children     []*Node
	Content      []byte
	LastModified *time.Time
}

type FileTree struct {
	Root *Node
}

func (ft *FileTree) Insert(file string) {
	ft.Root.Insert(file)
}

func (n *Node) insertHelper(newNode *Node) *Node {
	for _, child := range n.Children {
		if child.Name == newNode.Name {
			return child
		}
	}
	n.Children = append(n.Children, newNode)
	return newNode
}

func (n *Node) Insert(file string) {
	parts := strings.Split(file, "/")
	curr := n
	for i, part := range parts {
		isDir := i < len(parts)-1
		newNode := &Node{Name: part, IsDir: isDir, Parent: curr}
		curr = curr.insertHelper(newNode)
	}
}

func NewFileTree(input []string) *FileTree {
	root := &Node{Name: "", IsDir: true}
	root.Parent = root
	ft := &FileTree{Root: root}
	for _, file := range input {
		ft.Insert(file)
	}
	ft.Sort()
	return ft
}

func (ft *FileTree) Sort() {
	ft.Root.Sort()
}

func (n *Node) Sort() {
	sortByDir(n.Children)
	for _, child := range n.Children {
		child.Sort()
	}
}

func sortByDir(n []*Node) {
	sort.SliceStable(n, func(i, j int) bool {
		if n[i].IsDir == n[j].IsDir {
			return n[i].Name < n[j].Name
		}
		return n[i].IsDir
	})
}
