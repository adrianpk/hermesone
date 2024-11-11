package handler

import (
	"embed"
	"fmt"
	"os"
	"time"
)

const (
	defaultLayoutPath = "layout/default/default.html"
	backupLayoutPath  = "layout/default/default-%s.html.bak"
)

func InitDirs(dirs []string, layoutFS embed.FS) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(defaultLayoutPath); err == nil {
		timestamp := time.Now().Format("20060102150405")
		err := os.Rename(defaultLayoutPath, fmt.Sprintf(backupLayoutPath, timestamp))
		if err != nil {
			return err
		}
	}

	content, err := layoutFS.ReadFile(defaultLayoutPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(defaultLayoutPath, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
