package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrianpk/gohermes/internal/hermes"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewNewCmd creates a new content file with the specified properties.
//
// Example Usage:
//
// Assuming the hermes CLI tool is already built and available in your PATH, you can use the new command to create a new content file. Here is an example of how to use the command:
//
//	hermes new --name "How to Create New Content in Hermes" --type "article" --section "guides" --tags "tutorial,hermes" --author "Jane Doe,John Smith"
//
// Explanation::
//
// --name "How to Create New Content in Hermes": Specifies the name of the content. This will be converted to a slug (how-to-create-new-content-in-hermes.md).
// --type "article": Specifies the type of the content. This will determine the subdirectory under the section.
// --section "guides": Specifies the section of the content. This will determine the main directory under content.
// --tags "tutorial,hermes": Specifies the tags for the content.
// --author "Jane Doe,John Smith": Specifies the authors of the content.
//
// Resulting File Structure:
//
// After running the command, the following file structure will be created:
//
// content/
// └── guides/
//
//	└── article/
//	    └── how-to-create-new-content-in-hermes.md
//
// Resulting Markdown File:
//
// The content of how-to-create-new-content-in-hermes.md will be:
//
// ---
// title: "How to Create New Content in Hermes"
// description: ""
// date: "2023-10-01"
// published-at: "2023-10-01"
// last-modified: "2023-10-01"
// type: "article"
// section: "guides"
// slug: "how-to-create-new-content-in-hermes"
// author:
//   - "Jane Doe"
//   - "John Smith"
//
// tags:
//   - "tutorial"
//   - "hermes"
//
// ---
// # How to Create New Content in Hermes
func NewNewCmd() *cobra.Command {
	var name, contentType, section string
	var tags, authors []string

	cmd := &cobra.Command{
		Use:   "new",
		Short: "create a new content file",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Creating new content...")
			if name == "" {
				log.Println("Error: name is required")
				return
			}

			section :=  hermes.ValidSectionOrDef(section)
			contentType := hermes.ValidTypeOrDef(contentType)
			slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
			fileName := fmt.Sprintf("%s.md", slug)
			dirPath := filepath.Join("content", section, contentType)
			filePath := filepath.Join(dirPath, fileName)

			err := ensureDir(dirPath)
			if err != nil {
				log.Fatalf("Error: failed to create directories: %v", err)
			}

			now := time.Now().Format("2006-01-02")

			meta := hermes.Meta{
				Title:       name,
				Description: "",
				Date:        now,
				PublishedAt: now,
				CreatedAt:   now,
				UpdatedAt:   now,
				Type:        contentType,
				Section:     section,
				Slug:        slug,
				Authors:     authors,
				Tags:        tags,
				Draft:       true,
			}

			metaData, err := yaml.Marshal(meta)
			if err != nil {
				log.Fatalf("error: failed to marshal meta: %v", err)
			}

			content := fmt.Sprintf("---\n%s---\n# %s\n", string(metaData), name)
			err = os.WriteFile(filePath, []byte(content), 0644)
			if err != nil {
				log.Fatalf("error writing file: %v", err)
			}

			log.Printf("new content created at %s", filePath)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the content")
	cmd.Flags().StringVarP(&contentType, "type", "t", "", "type of the content")
	cmd.Flags().StringVarP(&section, "section", "s", "", "section of the content")
	cmd.Flags().StringSliceVarP(&tags, "tags", "g", []string{}, "tags for the content")
	cmd.Flags().StringSliceVarP(&authors, "author", "a", []string{}, "authors of the content")

	return cmd
}

// ensureDir creates the directory if it does not exist.
func ensureDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
