package handler

import (
	"bytes"
	"fmt"
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

const (
	pages    = "pages"
	articles = "articles"
	blog     = "blog"
	series   = "series"
)

const (
	noJekyllFile = ".nojekyll"
)

// GenHTML generates the HTML files from the markdown files.
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

				content, err := hermes.Parse(fileContent, path)
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
		return fmt.Errorf("error adding .nojekyll file: %w", err)
	}

	log.Println("content generated!")
	return nil
}

func determineOutputPath(relativePath string) string {
	parts := strings.Split(relativePath, osFileSep)

	if len(parts) < 2 {
		return outputPath(parts...)
	}

	section := parts[0]
	subdir := parts[1]
	needDir := needsCustomDir(subdir)

	switch section {
	case "root":
		if needDir {
			p := outputPath(append([]string{subdir}, parts[2:]...)...)
			return p
		} else {
			p := outputPath(append([]string{}, parts[2:]...)...)
			return p
		}

	case pages, articles:
		p := outputPath(parts[1:]...)
		return p

	default:
		if needDir {
			p := outputPath(append([]string{section, subdir}, parts[2:]...)...)
			return p
		} else {
			p := outputPath(append([]string{section}, parts[2:]...)...)
			return p
		}
	}
}

func needsCustomDir(dir string) bool {
	return dir != articles && dir != pages
}

func outputPath(parts ...string) string {
	trimmedPath := strings.TrimSuffix(strings.Join(parts, osFileSep), filepath.Ext(parts[len(parts)-1])) + ".html"
	return filepath.Join(outputRoot, trimmedPath)
}

// shouldRender checks if the markdown file is newer than the html file
// to determine if the html should be re-rendered.
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

// findLayout tries to find the layout file for the given markdown file.
func findLayout(path string) string {
	path = strings.TrimPrefix(path, "content/")

	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	dir := filepath.Dir(path)

	layoutPaths := []string{
		filepath.Join(defaultLayoutDir, dir, base+".html"),
		filepath.Join(defaultLayoutDir, dir, "default.html"),
		filepath.Join(defaultLayoutDir, filepath.Base(dir), "default.html"),
		filepath.Join(defaultLayoutDir, "default.html"),
	}

	for _, layoutPath := range layoutPaths {
		if _, err := os.Stat(layoutPath); err == nil {
			return layoutPath
		}
	}

	return ""
}

// copyImages copies the images from the markdown directory to the output directory.
func copyImages(mdPath, htmlPath string) error {
	rootPrefix := "root/"
	imageDir := strings.TrimSuffix(mdPath, filepath.Ext(mdPath))
	relativeImageDir := strings.TrimPrefix(imageDir, contentRoot+"/")

	if strings.HasPrefix(relativeImageDir, rootPrefix) {
		relativeImageDir = strings.TrimPrefix(relativeImageDir, rootPrefix)

		parts := strings.Split(relativeImageDir, osFileSep)

		if len(parts) == 1 {
			relativeImageDir = filepath.Join(imgDir, parts[0])
		} else if len(parts) > 1 {
			switch parts[0] {
			case blog, series:
				relativeImageDir = filepath.Join(imgDir, parts[0], parts[1])
			case articles, pages:
				relativeImageDir = filepath.Join(imgDir, parts[1])
			default:
				relativeImageDir = filepath.Join(imgDir, strings.Join(parts, osFileSep))
			}
		}
	} else {
		parts := strings.Split(relativeImageDir, osFileSep)

		if len(parts) > 2 && (parts[1] == articles || parts[1] == pages) {
			relativeImageDir = filepath.Join(imgDir, parts[0], parts[2])
		} else if len(parts) > 2 && (parts[1] == blog || parts[1] == series) {
			relativeImageDir = filepath.Join(imgDir, parts[0], parts[1], parts[2])
		} else {
			relativeImageDir = filepath.Join(imgDir, strings.Join(parts, osFileSep))
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

// addNoJekyll adds a .nojekyll file to the output directory.
func addNoJekyll() error {
	noJekyllPath := filepath.Join(outputRoot, noJekyllFile)
	if _, err := os.Stat(noJekyllPath); os.IsNotExist(err) {
		file, err := os.Create(noJekyllPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}
