package utils

import (
	"fmt"
	"os"
)

func CreateTempFile(data string, extension string) (*os.File, error) {
	file, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("*%s", extension))
	if err != nil {
		return nil, fmt.Errorf("unable to create new file: %w", err)
	}
	// TODO: find a better way?
	if data != "" {
		if err := os.WriteFile(file.Name(), []byte(data), 0o666); err != nil {
			return nil, fmt.Errorf("unable to write contents to file: %w", err)
		}
	}
	return file, nil
}
