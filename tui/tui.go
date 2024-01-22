package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/Wondrous27/s3-tui/bucket"
	"github.com/Wondrous27/s3-tui/object"
	"github.com/Wondrous27/s3-tui/tui/constants"
	tea "github.com/charmbracelet/bubbletea"
)

func StartTea(br bucket.S3Repository, or object.S3Repository) {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	constants.Br = &br
	constants.Or = &or

	m, _ := InitBuckets() // TODO: can we acknowledge this error
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := constants.P.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
