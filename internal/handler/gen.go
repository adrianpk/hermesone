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
	contentRoot = "content"
)

func GenHTML() error {
	err := filepath.Walk(contentRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			relativePath, err := filepath.Rel(contentRoot, path)
			if err != nil {
				return err
			}

			outputPath := determineOutputPath(relativePath)

			if shouldRender(path, outputPath) {
				fileContent, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				content, err := hermes.Parse(fileContent)
				if err != nil {
					return err
				}

				err = hermes.UpdateSection(path, &content.Meta)
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

					err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
					if err != nil {
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

	if err != nil {
		return err
	}

	err = addNoJekyll()
	if err != nil {
		log.Printf("error adding .nojekyll file: %v", err)
		return err
	}

	log.Println("content generated!")
	return nil
}

func determineOutputPath(relativePath string) string {
	parts := strings.Split(relativePath, string(os.PathSeparator))
	if len(parts) < 2 {
		return filepath.Join("output", relativePath[:len(relativePath)-3]+".html")
	}

	section := parts[0]
	subdir := parts[1]

	switch section {
	case "root":
		if subdir == "blog" || subdir == "series" {
			return filepath.Join("output", subdir, relativePath[len(section)+len(subdir)+2:len(relativePath)-3]+".html")
		}
		return filepath.Join("output", relativePath[len(section)+len(subdir)+2:len(relativePath)-3]+".html")
	default:
		if subdir == "blog" || subdir == "series" {
			return filepath.Join("output", section, subdir, relativePath[len(section)+len(subdir)+2:len(relativePath)-3]+".html")
		}
		return filepath.Join("output", section, relativePath[len(section)+len(subdir)+2:len(relativePath)-3]+".html")
	}
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

func addNoJekyll() error {
	noJekyllPath := filepath.Join("output", ".nojekyll")
	if _, err := os.Stat(noJekyllPath); os.IsNotExist(err) {
		file, err := os.Create(noJekyllPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}
