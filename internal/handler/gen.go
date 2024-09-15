package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianpk/gohermes/internal/gen"
)

func GenHTML() error {
	return filepath.Walk("content/root", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			content, err := gen.Process(fileContent)
			if err != nil {
				return err
			}

			layoutPath := findLayout(path)
			if layoutPath != "" {
				layout, err := os.ReadFile(layoutPath)
				if err != nil {
					return err
				}

				outputHTML := string(layout)

				// Replace metadata
				outputHTML = strings.Replace(outputHTML, "{{title}}", content.Meta.Title, -1)
				outputHTML = strings.Replace(outputHTML, "{{description}}", content.Meta.Description, -1)
				outputHTML = strings.Replace(outputHTML, "{{summary}}", content.Meta.Summary, -1)
				outputHTML = strings.Replace(outputHTML, "{{date}}", content.Meta.Date, -1)
				outputHTML = strings.Replace(outputHTML, "{{publish-date}}", content.Meta.PublishDate, -1)
				outputHTML = strings.Replace(outputHTML, "{{last-modified}}", content.Meta.LastModified, -1)
				outputHTML = strings.Replace(outputHTML, "{{type}}", content.Meta.Type, -1)
				outputHTML = strings.Replace(outputHTML, "{{section}}", content.Meta.Section, -1)
				outputHTML = strings.Replace(outputHTML, "{{slug}}", content.Meta.Slug, -1)
				outputHTML = strings.Replace(outputHTML, "{{image}}", content.Meta.Image, -1)
				outputHTML = strings.Replace(outputHTML, "{{social-image}}", content.Meta.SocialImage, -1)
				outputHTML = strings.Replace(outputHTML, "{{layout}}", content.Meta.Layout, -1)
				outputHTML = strings.Replace(outputHTML, "{{canonical-url}}", content.Meta.CanonicalURL, -1)
				outputHTML = strings.Replace(outputHTML, "{{locale}}", content.Meta.Locale, -1)
				outputHTML = strings.Replace(outputHTML, "{{robots}}", content.Meta.Robots, -1)
				outputHTML = strings.Replace(outputHTML, "{{excerpt}}", content.Meta.Excerpt, -1)
				outputHTML = strings.Replace(outputHTML, "{{permalink}}", content.Meta.Permalink, -1)

				outputHTML = strings.Replace(outputHTML, "{{draft}}", fmt.Sprintf("%v", content.Meta.Draft), -1)
				outputHTML = strings.Replace(outputHTML, "{{table-of-contents}}", fmt.Sprintf("%v", content.Meta.TableOfContents), -1)
				outputHTML = strings.Replace(outputHTML, "{{share}}", fmt.Sprintf("%v", content.Meta.Share), -1)
				outputHTML = strings.Replace(outputHTML, "{{featured}}", fmt.Sprintf("%v", content.Meta.Featured), -1)
				outputHTML = strings.Replace(outputHTML, "{{comments}}", fmt.Sprintf("%v", content.Meta.Comments), -1)

				outputHTML = strings.Replace(outputHTML, "{{authors}}", strings.Join(content.Meta.Authors, ", "), -1)
				outputHTML = strings.Replace(outputHTML, "{{categories}}", strings.Join(content.Meta.Categories, ", "), -1)
				outputHTML = strings.Replace(outputHTML, "{{tags}}", strings.Join(content.Meta.Tags, ", "), -1)
				outputHTML = strings.Replace(outputHTML, "{{keywords}}", strings.Join(content.Meta.Keywords, ", "), -1)

				outputHTML = strings.Replace(outputHTML, "{{sitemap-priority}}", fmt.Sprintf("%.2f", content.Meta.Sitemap.Priority), -1)
				outputHTML = strings.Replace(outputHTML, "{{sitemap-changefreq}}", content.Meta.Sitemap.ChangeFreq, -1)

				outputHTML = strings.Replace(outputHTML, "{{content}}", string(content.HTML), 1)

				html := []byte(outputHTML)
				outputPath := filepath.Join("output", filepath.Base(path[:len(path)-3]+".html"))
				err = os.WriteFile(outputPath, html, 0644)
				if err != nil {
					return err
				}

				fmt.Printf("Generated HTML for: %s\n", outputPath)
			} else {
				fmt.Println("No layout found for:", path)
			}
		}

		return nil
	})
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
