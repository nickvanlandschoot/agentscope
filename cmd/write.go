package cmd

import (
	"os"
)

func Write(path string, content string) error {
	b := []byte(content)
	return os.WriteFile(path, b, os.FileMode(0664))
}
