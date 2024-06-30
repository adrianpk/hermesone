package handler

import (
	"embed"
	"os"
	"time"
)

func InitDirs(dirs []string, layoutFS embed.FS) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat("layout/default/default.html"); err == nil {
		timestamp := time.Now().Format("20060102150405")
		err := os.Rename("layout/default/default.html", "layout/default/default-"+timestamp+".html.bak")
		if err != nil {
			return err
		}
	}

	content, err := layoutFS.ReadFile("layout/default/default.html")
	if err != nil {
		return err
	}

	err = os.WriteFile("layout/default/default.html", content, 0644)
	if err != nil {
		return err
	}

	return nil
}
