package handler

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianpk/gohermes/internal/hermes"
)

const (
	contentRoot = "content/root"
)

func GenHTML() error {
	return filepath.Walk(contentRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			outputPath := filepath.Join("output", filepath.Base(path[:len(path)-3]+".html"))

			if shouldRender(path, outputPath) {
				fileContent, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				content, err := hermes.Parse(fileContent)
				if err != nil {
					return err
				}

				layoutPath := findLayout(path)
				if layoutPath != "" {
					tmpl, err := template.New("webpage").Funcs(template.FuncMap{
						"safeHTML": func(s string) template.HTML { return template.HTML(s) },
					}).ParseFiles(layoutPath)

					if err != nil {
						log.Printf("error parsing template files: %v\n", err)
						return err
					}

					var tmplBuf bytes.Buffer
					err = tmpl.Execute(&tmplBuf, content)
					if err != nil {
						log.Printf("error executing template: %v\n", err)
						return err
					}

					outputFile, err := os.Create(outputPath)
					if err != nil {
						return err
					}
					defer outputFile.Close()

					_, err = tmplBuf.WriteTo(outputFile)
					if err != nil {
						return err
					}

				} else {
					log.Println("no layout found for:", path)
				}
			}
		}

		return nil
	})
}

func shouldRender(mdPath, htmlPath string) bool {
	htmlInfo, err := os.Stat(htmlPath)
	if os.IsNotExist(err) {
		return true
	}

	markdownInfo, err := os.Stat(mdPath)
	if err != nil {
		return false
	}

	render := markdownInfo.ModTime().After(htmlInfo.ModTime())

	return render
}

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
			return layoutPath
		}
	}

	return ""
}
