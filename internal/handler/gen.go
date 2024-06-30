package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func findLayout(path string) string {
	path = strings.TrimPrefix(path, "content/root/")

	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	dir := filepath.Dir(path)
	layoutPaths := []string{
		"layout/default/" + dir + "/" + base + ".html",
		"layout/default/" + dir + ".html",
		"layout/default/default.html",
	}

	for _, layoutPath := range layoutPaths {
		if _, err := os.Stat(layoutPath); err == nil {
			fmt.Println("Found layout:", layoutPath)
			return layoutPath
		} else {
			fmt.Println("Layout not found:", layoutPath)
		}
	}

	return ""
}

func GenerateHTML() error {
	err := filepath.Walk("content/root", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			html := blackfriday.Run(content)

			layoutPath := findLayout(path)
			if layoutPath != "" {
				layout, err := os.ReadFile(layoutPath)
				if err != nil {
					return err
				}

				html = []byte(strings.Replace(string(layout), "{{content}}", string(html), 1))
			} else {
				fmt.Println("No layout found for:", path)
			}

			outputPath := "output/" + filepath.Base(path)
			outputPath = outputPath[:len(outputPath)-3] + ".html"

			err = os.WriteFile(outputPath, html, 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
