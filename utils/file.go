package utils

import (
	"fmt"
	"os"
)

func CreateTempFile(data string) (*os.File, error) {
	file, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		return nil, fmt.Errorf("unable to create new file: %w", err)
	}

	if err := os.WriteFile(file.Name(), []byte(data), 0o666); err != nil {
		return nil, fmt.Errorf("unable to write contents to file: %w", err)
	}
	return file, nil
}
