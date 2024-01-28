package tui

import (
	"os"
	"os/exec"

	"github.com/Wondrous27/s3-tui/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func openEditorCmd(data string) tea.Cmd {
	file, err := utils.CreateTempFile(data)
	if err != nil {
		return func() tea.Msg {
			return errMsg{error: err}
		}
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, file.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err, file}
	})
}
