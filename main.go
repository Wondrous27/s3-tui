package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	CreationDate time.Time
	Name         string
	Objects      []*types.Object
}

func (i item) FilterValue() string { return i.Name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %v", index+1, i.Name)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
	client   *s3.Client
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "D":
			i := m.list.Index()
			m.list.RemoveItem(i)
			return m, nil
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.Name
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// if m.choice != "" {
	// 	quitMsg := fmt.Sprintf("You want to view %s, that's very cool this bucket can hold so much", m.choice)
	// 	return quitTextStyle.Render(quitMsg)
	// }
	if m.quitting {
		quitMsg := "you don't wanna check out any buckets, cool maaan"
		return quitTextStyle.Render(quitMsg)
	}
	return "\n" + m.list.View()
}

func listBuckets(client *s3.Client) []types.Bucket {
	out, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return out.Buckets
}

func listObjects(client *s3.Client, bucket string) []types.Object {
	out, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{Bucket: &bucket})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return out.Contents
}

func main() {
	/* TODO: Take region as a cmdline argument */
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		fmt.Println(err)
	}
	client := s3.NewFromConfig(cfg)
	buckets := listBuckets(client)

	var items []list.Item
	for _, b := range buckets {
		item := item{CreationDate: *b.CreationDate, Name: *b.Name, Objects: nil}
		items = append(items, item)
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Please choose a bucket"
	l.SetShowStatusBar(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, client: client}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
