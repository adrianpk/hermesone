package handler

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	imgDir        = "img"
	defaultLayout = "default.html"
)

var (
	defaultLayoutDir  = filepath.Join("layout", "default")
	defaultLayoutFile = filepath.Join(defaultLayoutDir, defaultLayout)
	backupLayoutPath  = filepath.Join(defaultLayoutDir, "default-%s.html.bak")
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

	if _, err := os.Stat(defaultLayoutFile); err == nil {
		timestamp := time.Now().Format("20060102150405")
		err := os.Rename(defaultLayoutFile, fmt.Sprintf(backupLayoutPath, timestamp))
		if err != nil {
			return err
		}
	}

	content, err := layoutFS.ReadFile(defaultLayoutFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(defaultLayoutFile, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
