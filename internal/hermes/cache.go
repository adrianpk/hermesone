package hermes

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

type FileData struct {
	Meta      Meta
	Path      string
	FileName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Cache struct {
	Data      []FileData
	CreatedAt time.Time
}

func BuildCache(contentDir string) (*Cache, error) {
	var metadata []FileData
	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".md" {
			meta, err := processFile(path)
			if err != nil {
				return err
			}
			metadata = append(metadata, meta)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	cache := &Cache{
		Data:      metadata,
		CreatedAt: time.Now(),
	}
	return cache, nil
}

func processFile(path string) (FileData, error) {
	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return FileData{}, err
	}

	parsedContent, err := parseFile(content, path)
	if err != nil {
		return FileData{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return FileData{}, err
	}

	createdAt := info.ModTime()
	updatedAt := info.ModTime()

	return FileData{
		Meta:      parsedContent.Meta,
		Path:      path,
		FileName:  filepath.Base(path),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// parseFile and Parse are similar.
// Integrate both at some point.
func parseFile(content []byte, path string) (Content, error) {
	var parsed Content

	if bytes.HasPrefix(content, separator) {
		parts := bytes.SplitN(content, separator, 3)
		if len(parts) == 3 {
			err := yaml.Unmarshal(parts[1], &parsed.Meta)
			if err != nil {
				return parsed, err
			}
			parsed.Markdown = parts[2]

		} else {
			parsed.Markdown = content
		}
	} else {
		parsed.Markdown = content
	}

	md := updateImgPaths(parsed.Markdown, path)

	parsed.HTML = string(blackfriday.Run(md))

	return parsed, nil
}

// Debug prints the cache content in a pretty, hierarchical way
func (c *Cache) Debug() {
	prettyJSON, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		fmt.Println("Failed to generate debug output:", err)
		return
	}
	fmt.Println(string(prettyJSON))
}
