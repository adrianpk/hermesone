package handler

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianpk/gohermes/internal/hermes"
)

const (
	contentRoot = "content"
	outputRoot  = "output"
)

func GenHTML() error {
	err := filepath.Walk(contentRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			relativePath, err := filepath.Rel(contentRoot, path)
			if err != nil {
				log.Printf("error getting relative path for %s: %v", path, err)
				return nil
			}

			outputPath := determineOutputPath(relativePath)

			if shouldRender(path, outputPath) {
				fileContent, err := os.ReadFile(path)
				if err != nil {
					log.Printf("error reading file %s: %v", path, err)
					return nil
				}

				content, err := hermes.Parse(fileContent)
				if err != nil {
					log.Printf("error parsing file %s: %v", path, err)
					return nil
				}

				err = hermes.UpdateSection(path, &content.Meta)
				if err != nil {
					log.Printf("error updating section for file %s: %v", path, err)
					return nil
				}

				layoutPath := findLayout(path)
				if layoutPath != "" {
					tmpl, err := template.New("webpage").Funcs(template.FuncMap{
						"safeHTML": func(s string) template.HTML { return template.HTML(s) },
					}).ParseFiles(layoutPath)

					if err != nil {
						log.Printf("error parsing template files for %s: %v", path, err)
						return nil
					}

					var tmplBuf bytes.Buffer
					err = tmpl.Execute(&tmplBuf, content)
					if err != nil {
						log.Printf("error executing template for %s: %v", path, err)
						return nil
					}

					err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
					if err != nil {
						log.Printf("error creating directories for %s: %v", outputPath, err)
						return nil
					}

					outputFile, err := os.Create(outputPath)
					if err != nil {
						log.Printf("error creating output file %s: %v", outputPath, err)
						return nil
					}
					defer outputFile.Close()

					_, err = tmplBuf.WriteTo(outputFile)
					if err != nil {
						log.Printf("error writing to output file %s: %v", outputPath, err)
						return nil
					}

					log.Printf("### copied %s to %s", path, outputPath)
					err = copyImages(path, outputPath)
					if err != nil {
						log.Printf("error copying images for %s: %v", path, err)
						return nil
					}

				} else {
					log.Printf("no layout found for %s", path)
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
		return filepath.Join(outputRoot, relativePath[:len(relativePath)-3]+".html")
	}

	section := parts[0]
	subdir := parts[1]

	switch section {
	case "root":
		if subdir == "blog" || subdir == "series" {
			return filepath.Join(outputRoot, subdir, strings.Join(parts[2:], string(os.PathSeparator))[:len(relativePath)-len(section)-len(subdir)-2]+".html")
		}
		if subdir == "articles" || subdir == "pages" {
			return filepath.Join(outputRoot, strings.Join(parts[2:], string(os.PathSeparator))[:len(relativePath)-len(section)-len(subdir)-2]+".html")
		}
		return filepath.Join(outputRoot, strings.Join(parts[1:], string(os.PathSeparator))[:len(relativePath)-len(section)-1]+".html")
	case "pages", "articles":
		return filepath.Join(outputRoot, strings.Join(parts[1:], string(os.PathSeparator))[:len(relativePath)-len(section)-1]+".html")
	default:
		if subdir == "blog" || subdir == "series" {
			return filepath.Join(outputRoot, section, subdir, strings.Join(parts[2:], string(os.PathSeparator))[:len(relativePath)-len(section)-len(subdir)-2]+".html")
		}
		if subdir == "articles" || subdir == "pages" {
			return filepath.Join(outputRoot, section, strings.Join(parts[2:], string(os.PathSeparator))[:len(relativePath)-len(section)-len(subdir)-2]+".html")
		}
		return filepath.Join(outputRoot, section, strings.Join(parts[1:], string(os.PathSeparator))[:len(relativePath)-len(section)-1]+".html")
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

func copyImages(mdPath, htmlPath string) error {
	imageDir := strings.TrimSuffix(mdPath, filepath.Ext(mdPath))
	relativeImageDir := strings.TrimPrefix(imageDir, contentRoot+"/")

	if strings.HasPrefix(relativeImageDir, "root/") {
		relativeImageDir = strings.TrimPrefix(relativeImageDir, "root/")

		parts := strings.Split(relativeImageDir, string(os.PathSeparator))

		if len(parts) == 1 {
			// .md file is directly under content/root (e.g., content/root/index.md)
			relativeImageDir = filepath.Join("img", parts[0])
		} else if len(parts) > 1 {
			switch parts[0] {
			case "blog", "series":
				relativeImageDir = filepath.Join("img", parts[0], parts[1])
			case "articles", "pages":
				// For images in content/root/articles or content/root/pages, copy to output/img/{content-name-without-ext}
				relativeImageDir = filepath.Join("img", parts[1])
			default:
				relativeImageDir = filepath.Join("img", strings.Join(parts, string(os.PathSeparator)))
			}
		}
	} else {
		// Sections
		parts := strings.Split(relativeImageDir, string(os.PathSeparator))

		if len(parts) > 2 && (parts[1] == "articles" || parts[1] == "pages") {
			relativeImageDir = filepath.Join("img", parts[0], parts[2])
		} else if len(parts) > 2 && (parts[1] == "blog" || parts[1] == "series") {
			relativeImageDir = filepath.Join("img", parts[0], parts[1], parts[2])
		} else {
			relativeImageDir = filepath.Join("img", strings.Join(parts, string(os.PathSeparator)))
		}
	}

	outputImageDir := filepath.Join(outputRoot, relativeImageDir)

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relativePath, err := filepath.Rel(imageDir, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(outputImageDir, relativePath)

			err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			if err != nil {
				return err
			}

			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func addNoJekyll() error {
	noJekyllPath := filepath.Join(outputRoot, ".nojekyll")
	if _, err := os.Stat(noJekyllPath); os.IsNotExist(err) {
		file, err := os.Create(noJekyllPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}
