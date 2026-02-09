package cmd

import (
	"os"
)

func Create(fileName string) error {
	file, err := os.Create(fileName)
	defer file.Close()
	return err
}
