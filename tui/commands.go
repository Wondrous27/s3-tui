package tui

import (
	"os"
	"os/exec"

	"github.com/Wondrous27/s3-tui/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func openEditorCmd(data string, extension string) tea.Cmd {
	file, err := utils.CreateTempFile(data, extension)
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

// func (m Object) updateObjectCmd(fileName string) tea.Cmd {
// 	return func() tea.Msg {
// 		file, _ := os.Open(fileName)
// 		key := m.object.Key
// 		bucket := m.activeBucketName
// 		err := constants.Or.PutObject(file, bucket, key)
// 		if err != nil {
// 			return errMsg{fmt.Errorf("cannot read file in createEntryCmd: %v", err)}
// 		}
// 		return m.setupObject()
// 	}
// }
