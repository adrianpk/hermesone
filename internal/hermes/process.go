package hermes

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

func Parse(content []byte) (Content, error) {
	var parsed Content

	if bytes.HasPrefix(content, []byte("---\n")) {
		parts := bytes.SplitN(content, []byte("---\n"), 3)
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

	parsed.HTML = string(blackfriday.Run(parsed.Markdown))

	return parsed, nil
}

// UpdateSection updates the section property in the Meta struct based on the directory of the markdown file.
func UpdateSection(mdPath string, meta *Meta) error {
	dir := filepath.Dir(mdPath)
	dir = strings.TrimPrefix(dir, "content/")
	parts := strings.Split(dir, string(os.PathSeparator))
	if len(parts) < 2 {
		return fmt.Errorf("invalid directory structure for file %s", mdPath)
	}
	section := parts[0]

	if meta.Section != section {
		meta.Section = section

		fileContent, err := os.ReadFile(mdPath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", mdPath, err)
		}

		parts := bytes.SplitN(fileContent, []byte("\n---\n"), 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid file format")
		}

		updatedFrontMatter, err := yaml.Marshal(meta)
		if err != nil {
			return fmt.Errorf("error marshaling updated front matter: %v", err)
		}

		updatedContent := append([]byte("---\n"), updatedFrontMatter...)
		updatedContent = append(updatedContent, []byte("---\n")...)
		updatedContent = append(updatedContent, parts[1]...)

		err = os.WriteFile(mdPath, updatedContent, 0644)
		if err != nil {
			return fmt.Errorf("error writing updated content to file %s: %v", mdPath, err)
		}
	}

	return nil
}
