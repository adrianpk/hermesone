package hermes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	Root          string
	Data          map[string]FileData
	All           []FileData
	BySection     map[string][]FileData
	BySectionType map[string]map[string][]FileData
	ByTags        map[string][]FileData
	ByPath        map[string]FileData
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

			pp.Data[relativePath] = fileData
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// FindFileData finds the file data by the relative path.
func (pp *PreProcessor) FindFileData(relativePath string) (FileData, bool) {
	fileData, exists := pp.Data[relativePath]
	return fileData, exists
}

// GetPublishedContent returns a list of published content ordered by PublishedAt in descending order
func (pp *PreProcessor) GetPublishedContent() []FileData {
	var publishedContent []FileData
	for _, fileData := range pp.Data {
		if fileData.Published {
			publishedContent = append(publishedContent, fileData)
		}
	}

	sort.Slice(publishedContent, func(i, j int) bool {
		publishedAtI, _ := time.Parse(time.RFC3339, publishedContent[i].Meta.PublishedAt)
		publishedAtJ, _ := time.Parse(time.RFC3339, publishedContent[j].Meta.PublishedAt)
		return publishedAtI.After(publishedAtJ)
	})

	return publishedContent
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

	pp.All = []FileData{}
	pp.BySection = make(map[string][]FileData)
	pp.BySectionType = make(map[string]map[string][]FileData)
	pp.ByTags = make(map[string][]FileData)
	pp.ByPath = make(map[string]FileData)

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
		sectionUpdated := pp.updateSection(filePath, &fileData.Meta)
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

		if fileData.Published {
			pp.All = append(pp.All, fileData)

			section := fileData.Meta.Section
			pp.BySection[section] = append(pp.BySection[section], fileData)

			sectionType := fileData.Meta.Type
			if pp.BySectionType[section] == nil {
				pp.BySectionType[section] = make(map[string][]FileData)
			}
			pp.BySectionType[section][sectionType] = append(pp.BySectionType[section][sectionType], fileData)

			for _, tag := range fileData.Meta.Tags {
				pp.ByTags[tag] = append(pp.ByTags[tag], fileData)
			}

			pp.ByPath[path] = fileData
		}
	}

	pp.sortAll()
	pp.sortBySection()
	pp.sortBySectionType()
	pp.sortByTags()

	return nil
}

// updateSection updates the section property in the Meta struct based on the directory of the markdown file.
// Returns true if the section was updated, false otherwise.
func (pp *PreProcessor) updateSection(mdPath string, meta *Meta) bool {
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

// GetAllPublished returns all published content ordered by PublishedAt date (newest first)
func (pp *PreProcessor) GetAllPublished() []FileData {
	return pp.All
}

// GetPublishedBySection returns published content for a specific section ordered by PublishedAt date (newest first)
func (pp *PreProcessor) GetPublishedBySection(section string) []FileData {
	return pp.BySection[section]
}

// GetPublishedBySectionType returns published content for a specific section and type ordered by PublishedAt date (newest first)
func (pp *PreProcessor) GetPublishedBySectionType(section, sectionType string) []FileData {
	return pp.BySectionType[section][sectionType]
}

// GetPublishedByTag returns published content for a specific tag ordered by PublishedAt date (newest first)
func (pp *PreProcessor) GetPublishedByTag(tag string) []FileData {
	return pp.ByTags[tag]
}

// GetPublishedByPath returns published content for a specific path
func (pp *PreProcessor) GetPublishedByPath(path string) (FileData, bool) {
	fileData, exists := pp.ByPath[path]
	return fileData, exists
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

func (pp *PreProcessor) sortAll() {
	sort.Slice(pp.All, func(i, j int) bool {
		publishedAtI, _ := time.Parse(time.RFC3339, pp.All[i].Meta.PublishedAt)
		publishedAtJ, _ := time.Parse(time.RFC3339, pp.All[j].Meta.PublishedAt)
		return publishedAtI.After(publishedAtJ)
	})
}

func (pp *PreProcessor) sortBySection() {
	for section, files := range pp.BySection {
		sort.Slice(files, func(i, j int) bool {
			publishedAtI, _ := time.Parse(time.RFC3339, files[i].Meta.PublishedAt)
			publishedAtJ, _ := time.Parse(time.RFC3339, files[j].Meta.PublishedAt)
			return publishedAtI.After(publishedAtJ)
		})
		pp.BySection[section] = files
	}
}

func (pp *PreProcessor) sortBySectionType() {
	for section, types := range pp.BySectionType {
		for sectionType, files := range types {
			sort.Slice(files, func(i, j int) bool {
				publishedAtI, _ := time.Parse(time.RFC3339, files[i].Meta.PublishedAt)
				publishedAtJ, _ := time.Parse(time.RFC3339, files[j].Meta.PublishedAt)
				return publishedAtI.After(publishedAtJ)
			})
			pp.BySectionType[section][sectionType] = files
		}
	}
}

func (pp *PreProcessor) sortByTags() {
	for tag, files := range pp.ByTags {
		sort.Slice(files, func(i, j int) bool {
			publishedAtI, _ := time.Parse(time.RFC3339, files[i].Meta.PublishedAt)
			publishedAtJ, _ := time.Parse(time.RFC3339, files[j].Meta.PublishedAt)
			return publishedAtI.After(publishedAtJ)
		})
		pp.ByTags[tag] = files
	}
}

// GetAllPublishedPaginated returns a paginated list of all published content
func (pp *PreProcessor) GetAllPublishedPaginated(page, pageSize int) ([]FileData, error) {
	return Paginate(pp.All, page, pageSize)
}

// GetPublishedBySectionPaginated returns a paginated list of published content for a specific section
func (pp *PreProcessor) GetPublishedBySectionPaginated(section string, page, pageSize int) ([]FileData, error) {
	return Paginate(pp.BySection[section], page, pageSize)
}

// GetPublishedBySectionTypePaginated returns a paginated list of published content for a specific section and type
func (pp *PreProcessor) GetPublishedBySectionTypePaginated(section, sectionType string, page, pageSize int) ([]FileData, error) {
	return Paginate(pp.BySectionType[section][sectionType], page, pageSize)
}

// GetPublishedByTagPaginated returns a paginated list of published content for a specific tag
func (pp *PreProcessor) GetPublishedByTagPaginated(tag string, page, pageSize int) ([]FileData, error) {
	return Paginate(pp.ByTags[tag], page, pageSize)
}
