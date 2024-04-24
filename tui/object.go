package tui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Wondrous27/s3-tui/object"
	"github.com/Wondrous27/s3-tui/tui/constants"
	"github.com/alecthomas/chroma/lexers"
	"github.com/muesli/termenv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	gansi "github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/lipgloss"
)

type errMsg struct{ error }

type UpdatedObject *object.Object

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
	object           object.Object
	quitting         bool
}

// Init run any intial IO on program start
func (m Object) Init() tea.Cmd {
	return nil
}

// initialize the objectui model for your program
func InitObject(bucketName, key string) (tea.Model, tea.Cmd) {
	m := Object{activeBucketName: bucketName}
	top, right, bottom, left := constants.DocStyle.GetMargin()
	m.viewport = viewport.New(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	m.viewport.Style = lipgloss.NewStyle().Align(lipgloss.Bottom)

	obj, ok := m.setupObject(bucketName, key).(UpdatedObject)
	if !ok {
		log.Println("failed to setup object")
		return m, tea.Quit
	}
	m.object = *obj
	m.setViewportContent()
	return &m, nil
}

func (m *Object) setupObject(bucketName, key string) tea.Msg {
	obj, err := constants.Or.GetObject(bucketName, key)
	if err != nil {
		return errMsg{fmt.Errorf("cannot get content: %v", err)}
	}
	return UpdatedObject(obj)
}

func (m *Object) setViewportContent() {
	var str string
	var err error
	content := object.FormatObject(m.object)
	if m.isSelectedMarkdown() {
		str, err = glamour.Render(content, "dark")
		if err != nil {
			m.error = "could not render content with glamour"
		}
	} else {
		str, err = m.renderFile(m.object.Key, content)
		log.Println("rendering file", str)
		if err != nil {
			log.Println("error rendering file", err)
			m.error = "could not render content with renderFile"
		}

	}
	m.viewport.SetContent(str)
}

func (m *Object) renderFile(path, content string) (string, error) {
	lexer := lexers.Match(path)
	if path == "" {
		lexer = lexers.Analyse(content)
	}
	lang := ""
	if lexer != nil && lexer.Config() != nil {
		lang = lexer.Config().Name
	}

	formatter := &gansi.CodeBlockElement{
		Code:     content,
		Language: lang,
	}
	s := strings.Builder{}
	st := constants.StyleConfig()
	var u uint
	st.CodeBlock.Margin = &u
	rc := gansi.NewRenderContext(gansi.Options{
		ColorProfile: termenv.TrueColor,
		Styles:       st,
	})
	err := formatter.Render(&s, rc)
	if err != nil {
		return "", err
	}

	return s.String(), nil
}

func (m Object) isSelectedMarkdown() bool {
	var lang string
	lexer := lexers.Match(m.object.Key)
	if lexer == nil {
		lexer = lexers.Analyse(m.object.Content)
	}
	if lexer != nil && lexer.Config() != nil {
		lang = lexer.Config().Name
	}
	return lang == "markdown"
}

func (m Object) helpView() string {
	// TODO: use the keymaps to populate the help string
	return constants.HelpStyle(
		"\n ↑/↓: h/j/k/l navigate  • esc: back • e: edit object • d: delete object • q: quit\n",
	)
}

// View return the text UI to be output to the terminal
func (m Object) View() string {
	if m.quitting {
		return ""
	}

	formatted := lipgloss.JoinVertical(
		lipgloss.Left,
		"\n",
		m.viewport.View(),
		m.helpView(),
		constants.ErrStyle(m.error),
	)
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

	case UpdatedObject:
		m.object = *msg

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Edit):
			fileContent := m.object.Content
			keys := strings.Split(m.object.Key, "/")
			fileName := keys[len(keys)-1]
			return m, openEditorCmd(fileContent, fileName)
		case key.Matches(msg, constants.Keymap.Create):
			return m, nil
		case key.Matches(msg, constants.Keymap.Back):
			return InitBuckets()

		case key.Matches(msg, constants.Keymap.Prev):
			tree := InitTree(m.activeBucketName)
			return tree.Update(constants.WindowSize)

		case key.Matches(msg, constants.Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.setViewportContent()
	return m, tea.Batch(cmds...)
}
