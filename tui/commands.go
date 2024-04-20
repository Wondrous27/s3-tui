package tui

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Wondrous27/s3-tui/tui/constants"
	"github.com/Wondrous27/s3-tui/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func openEditorCmd(data, extension string) tea.Cmd {
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

func (m Object) updateObjectCmd(fileName string) tea.Cmd {
	return func() tea.Msg {
		file, _ := os.Open(fileName)
		key := m.object.Key
		bucket := m.activeBucketName
		err := constants.Or.PutObject(file, bucket, key)
		if err != nil {
			return errMsg{fmt.Errorf("[updateObjectCmd]: cannot put object %v", err)}
		}
		return m.setupObject(bucket, key)
	}
}

func (f Tree) createObjectCommand(fileName, s3Key string) tea.Cmd {
	return func() tea.Msg {
		file, _ := os.Open(fileName)
		bucket := f.BucketName
		err := constants.Or.PutObject(file, bucket, s3Key)
		log.Printf("putting object with fileName %s, bucket %s, key %s", fileName, bucket, s3Key)
		if err != nil {
			return errMsg{fmt.Errorf("[createObjectCommand] cannot put object %v", err)}
		}
		return f.setupTree(f.BucketName)
	}
}
