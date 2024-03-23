package tui

import (
	"fmt"
	"os"

	"github.com/Wondrous27/s3-tui/bucket"
	"github.com/Wondrous27/s3-tui/object"
	"github.com/Wondrous27/s3-tui/tui/constants"
	tea "github.com/charmbracelet/bubbletea"
)

func StartTea(br *bucket.S3Repository, or *object.S3Repository) {
	constants.Br = br
	constants.Or = or

	m, err := InitBuckets()
	if err != nil {
		fmt.Println("Error initializing buckets", err)
		os.Exit(1)
	}

	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := constants.P.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
