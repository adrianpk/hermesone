package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianpk/gohermes/internal/fm"
	"github.com/russross/blackfriday/v2"
)

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

			parsed, err := fm.Preprocess(content)
			if err != nil {
				return err
			}

			html := blackfriday.Run(parsed.Markdown)

			layoutPath := findLayout(path)
			if layoutPath != "" {
				layout, err := os.ReadFile(layoutPath)
				if err != nil {
					return err
				}

				htmlStr := strings.Replace(string(layout), "{{content}}", string(html), 1)

				meta := parsed.Metadata
				htmlStr = strings.Replace(htmlStr, "{{title}}", meta.Title, 1)
				htmlStr = strings.Replace(htmlStr, "{{description}}", meta.Description, 1)
				htmlStr = strings.Replace(htmlStr, "{{date}}", meta.Date, 1)
				htmlStr = strings.Replace(htmlStr, "{{publish-date}}", meta.PublishDate, 1)
				htmlStr = strings.Replace(htmlStr, "{{last-modified}}", meta.LastModified, 1)
				htmlStr = strings.Replace(htmlStr, "{{type}}", meta.Type, 1)
				htmlStr = strings.Replace(htmlStr, "{{section}}", meta.Section, 1)
				htmlStr = strings.Replace(htmlStr, "{{slug}}", meta.Slug, 1)
				htmlStr = strings.Replace(htmlStr, "{{summary}}", meta.Summary, 1)
				htmlStr = strings.Replace(htmlStr, "{{image}}", meta.Image, 1)
				htmlStr = strings.Replace(htmlStr, "{{social-image}}", meta.SocialImage, 1)
				htmlStr = strings.Replace(htmlStr, "{{layout}}", meta.Layout, 1)
				htmlStr = strings.Replace(htmlStr, "{{canonical-url}}", meta.CanonicalURL, 1)
				htmlStr = strings.Replace(htmlStr, "{{locale}}", meta.Locale, 1)
				htmlStr = strings.Replace(htmlStr, "{{robots}}", meta.Robots, 1)
				htmlStr = strings.Replace(htmlStr, "{{excerpt}}", meta.Excerpt, 1)
				htmlStr = strings.Replace(htmlStr, "{{permalink}}", meta.Permalink, 1)

				htmlStr = strings.Replace(htmlStr, "{{draft}}", fmt.Sprintf("%v", meta.Draft), 1)
				htmlStr = strings.Replace(htmlStr, "{{table-of-contents}}", fmt.Sprintf("%v", meta.TableOfContents), 1)
				htmlStr = strings.Replace(htmlStr, "{{share}}", fmt.Sprintf("%v", meta.Share), 1)
				htmlStr = strings.Replace(htmlStr, "{{featured}}", fmt.Sprintf("%v", meta.Featured), 1)
				htmlStr = strings.Replace(htmlStr, "{{comments}}", fmt.Sprintf("%v", meta.Comments), 1)

				htmlStr = strings.Replace(htmlStr, "{{authors}}", strings.Join(meta.Authors, ", "), 1)
				htmlStr = strings.Replace(htmlStr, "{{categories}}", strings.Join(meta.Categories, ", "), 1)
				htmlStr = strings.Replace(htmlStr, "{{tags}}", strings.Join(meta.Tags, ", "), 1)
				htmlStr = strings.Replace(htmlStr, "{{keywords}}", strings.Join(meta.Keywords, ", "), 1)

				htmlStr = strings.Replace(htmlStr, "{{sitemap.priority}}", fmt.Sprintf("%.2f", meta.Sitemap.Priority), 1)
				htmlStr = strings.Replace(htmlStr, "{{sitemap.changefreq}}", meta.Sitemap.ChangeFreq, 1)

				html = []byte(htmlStr)
			} else {
				log.Println("No layout found for:", path)
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
			log.Println("found layout:", layoutPath)
			return layoutPath
		} else {
			log.Println("layout not found:", layoutPath)
		}
	}

	return ""
}
