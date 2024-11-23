package handler

import (
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/adrianpk/gohermes/internal/hermes"
)

var (
	osFileSep = string(os.PathSeparator)
)

func InitDirs(dirs []string, layoutFS embed.FS) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(hermes.DefLayoutFile); err == nil {
		timestamp := time.Now().Format("20060102150405")
		err := os.Rename(hermes.DefLayoutFile, fmt.Sprintf(hermes.BakLayoutPathFormat, timestamp))
		if err != nil {
			return err
		}
	}

	content, err := layoutFS.ReadFile(hermes.DefLayoutFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(hermes.DefLayoutFile, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
