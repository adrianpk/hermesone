package hermes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// FileData represents the data of a file
type FileData struct {
	Meta      Meta   `json:"meta"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Published bool
}

// PreProcessor adjusts file metadata and content.
type PreProcessor struct {
	Data map[string]FileData
	Root string
}

// NewPreProcessor creates a new PreProcessor instance
func NewPreProcessor(root string) *PreProcessor {
	return &PreProcessor{
		Data: make(map[string]FileData),
		Root: root,
	}
}

// Build builds the cache from the given root directory
func (pp *PreProcessor) Build() error {
	fmt.Printf("building cache: %s\n", pp.Root)
	err := filepath.Walk(pp.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//fmt.Printf("processing file: %s\n", path)

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("error reading file: %s\n", err)
				return err
			}

			parts := bytes.SplitN(fileContent, []byte("\n---\n"), 2)
			if len(parts) < 2 {
				fmt.Printf("invalid frontmatter format in file: %s\n", path)
				return fmt.Errorf("invalid frontmatter format in file: %s", path)
			}

			var meta Meta
			err = yaml.Unmarshal(parts[0], &meta)
			if err != nil {
				fmt.Printf("error unmarshalling frontmatter: %s\n", err)
				return err
			}

			fileData := FileData{
				Meta:      meta,
				Content:   string(parts[1]),
				CreatedAt: info.ModTime().Format(time.RFC3339),
				UpdatedAt: info.ModTime().Format(time.RFC3339),
			}

			fileData.UpdatePublishedStatus()

			relativePath, err := filepath.Rel(pp.Root, path)
			if err != nil {
				fmt.Printf("error getting relative path: %s\n", err)
				return err
			}

			//fmt.Printf("adding file to preprocessor cache: %s\n", relativePath)

			pp.Data[relativePath] = fileData
		}

		return nil
	})

	return err
}

// FindFileData finds the file data by the relative path.
func (pp *PreProcessor) FindFileData(relativePath string) (FileData, bool) {
	fileData, exists := pp.Data[relativePath]
	return fileData, exists
}

// Debug prints the cache content in a pretty, hierarchical way
func (pp *PreProcessor) Debug() {
	prettyJSON, err := json.MarshalIndent(pp.Data, "", "  ")
	if err != nil {
		fmt.Println("failed to generate debug output:", err)
		return
	}
	fmt.Println(string(prettyJSON))
}

// Sync aligns data obtained from the file with the Meta in FileData struct
func (pp *PreProcessor) Sync() error {
	const marginOfError = time.Second

	for path, fileData := range pp.Data {
		filePath := filepath.Join(pp.Root, path)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		updated := false

		if fileData.Meta.CreatedAt == "" {
			fileData.Meta.CreatedAt = fileInfo.ModTime().Format(time.RFC3339)
			updated = true
		}

		if fileData.Meta.UpdatedAt == "" {
			fileData.Meta.UpdatedAt = fileInfo.ModTime().Format(time.RFC3339)
			updated = true
		}

		fileModTime, err := time.Parse(time.RFC3339, fileData.Meta.UpdatedAt)
		if err != nil {
			return err
		}

		if fileModTime.Sub(fileInfo.ModTime()) > marginOfError || fileInfo.ModTime().Sub(fileModTime) > marginOfError {
			fileData.Meta.UpdatedAt = time.Now().Format(time.RFC3339)
			updated = true
		}

		// Update the section accordingly without writing the file again
		sectionUpdated := CorrectSection(filePath, &fileData.Meta)
		if sectionUpdated {
			updated = true
		}

		if updated {
			var frontmatter bytes.Buffer
			err = yaml.NewEncoder(&frontmatter).Encode(fileData.Meta)
			if err != nil {
				return err
			}

			newContent := append([]byte("---\n"), frontmatter.Bytes()...)
			newContent = append(newContent, []byte("---\n")...)
			newContent = append(newContent, []byte(fileData.Content)...)

			currentTime := time.Now().Format(time.RFC3339)
			fileData.Meta.UpdatedAt = currentTime

			err = os.WriteFile(filePath, newContent, fileInfo.Mode())
			if err != nil {
				return err
			}

			pp.Data[path] = fileData

			fmt.Printf("updated file: %s\n", filePath)
		}
	}

	return nil
}

func CorrectSection(mdPath string, meta *Meta) bool {
	dir := filepath.Dir(mdPath)
	dir = strings.TrimPrefix(dir, "content/")
	parts := strings.Split(dir, string(os.PathSeparator))
	if len(parts) < 2 {
		return false
	}
	section := parts[0]

	if meta.Section != section {
		meta.Section = section
		return true
	}

	return false
}

func (fd *FileData) UpdatePublishedStatus() {
	if fd.Meta.PublishedAt == "" {
		fd.Published = false
		return
	}

	publishedAt, err := time.Parse(time.RFC3339, fd.Meta.PublishedAt)
	if err != nil {
		fd.Published = false
		return
	}

	fd.Published = time.Now().After(publishedAt)
}
