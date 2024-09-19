package handler

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/adrianpk/gohermes/internal/hermes"
)

func Upgrade(dirs []string, layoutFS embed.FS) error {
	err := hermes.CheckHermes()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		log.Printf("creating directory: %s", dir)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	files := []string{
		"layout/default/default.html",
		"layout/default/article/default.html",
		"layout/default/blog/default.html",
		"layout/default/page/default.html",
		"layout/default/series/default.html",
	}

	for _, file := range files {
		log.Printf("processing file: %s", file)

		if _, err := os.Stat(file); err == nil {
			timestamp := time.Now().Format("060102150405")
			ext := filepath.Ext(file)
			newName := fmt.Sprintf("%s.%s%s", file[:len(file)-len(ext)], timestamp, ext)

			log.Printf("renaming existing file: %s to %s", file, newName)
			err = os.Rename(file, newName)
			if err != nil {
				return fmt.Errorf("failed to rename file %s to %s: %w", file, newName, err)
			}
		}

		content, err := layoutFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file, err)
		}

		log.Printf("writing file: %s", file)

		err = os.WriteFile(file, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", file, err)
		}
	}

	return nil
}
