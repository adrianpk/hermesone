package handler

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
)

func GenerateHTML() error {
	err := filepath.Walk("content/root", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			html := blackfriday.Run(content)

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
