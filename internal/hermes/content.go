package hermes

import (
	"errors"
	"log"
	"time"
)

var (
	validTypes = map[string]bool{
		ContentType.Article: true,
		ContentType.Blog:    true,
		ContentType.Series:  true,
	}

	indexTypes = map[string]bool{
		ContentType.Page: true,
	}

	indexableTypes = map[string]bool{
		ContentType.Article: true,
		ContentType.Blog:    true,
		ContentType.Series:  true,
	}
)

type Content struct {
	Meta     Meta
	Markdown []byte
	HTML     string
}

type Meta struct {
	Title           string   `yaml:"title"`        // Title of the content
	Description     string   `yaml:"description"`  // Description of the content
	FilePath        string   `json:"file_path"`    // File path of the content
	Summary         string   `yaml:"summary"`      // Summary of the content
	Date            string   `yaml:"date"`         // Date of the content
	PublishedAt     string   `yaml:"published-at"` // Publish date of the content
	CreatedAt       string   `yaml:"created-at"`   // Created date of the content
	UpdatedAt       string   `yaml:"updated-at"`   // Updated date of the content
	Type            string   `yaml:"type"`         // Type of the content
	Section         string   `yaml:"section"`      // Section of the content
	Slug            string   `yaml:"slug"`         // Slug of the content
	HeaderImage     string   `yaml:"header-image"` // HeaderImage of the content
	SocialImage     string   `yaml:"social-image"` // Social image of the content
	Layout          string   `yaml:"layout"`
	CanonicalURL    string   `yaml:"canonical-url"`
	Locale          string   `yaml:"locale"`  // Locale of the content
	Robots          string   `yaml:"robots"`  // Robots of the content
	Excerpt         string   `yaml:"excerpt"` // Excerpt of the content
	Permalink       string   `yaml:"permalink"`
	Draft           bool     `yaml:"draft"` // Draft of the content
	TableOfContents bool     `yaml:"table-of-contents"`
	Share           bool     `yaml:"share"`
	Featured        bool     `yaml:"featured"`
	Comments        bool     `yaml:"comments"`
	Authors         []string `yaml:"author"` // Authors of the content
	Categories      []string `yaml:"categories"`
	Tags            []string `yaml:"tags"` // Tags of the content
	Keywords        []string `yaml:"keywords"`
	Sitemap         Sitemap  `yaml:"sitemap"`
	Debug           bool     `yaml:"debug"`
}

type Sitemap struct {
	Priority   float64 `yaml:"priority"`
	ChangeFreq string  `yaml:"changefreq"`
}

func (m *Meta) PublicationDate() (t time.Time, err error) {
	if m.PublishedAt == "" {
		return t, errors.New("no publication date defined")
	}

	formats := []string{
		time.RFC3339,                // 2006-01-02T15:04:05Z07:00
		"2006-01-02",                // 2006-01-02
		"2006-01-02 15:04:05",       // 2006-01-02 15:04:05
		"2006-01-02 15:04",          // 2006-01-02 15:04
		"2006-01-02 15:04:05 -0700", // 2006-01-02 15:04:05 -0700
	}

	for _, format := range formats {
		publishedAt, err := time.Parse(format, m.PublishedAt)
		if err == nil {
			return publishedAt, nil
		}
	}

	log.Printf("error parsing publication date: %v", err)
	return t, errors.New("invalid publication date format")
}

func (m *Meta) IsPublished() bool {
	if m.Draft {
		return false
	}

	pd, err := m.PublicationDate()
	if err != nil {
		return false
	}

	return time.Now().After(pd)
}

func (m *Meta) IsIndexable() bool {
	if !m.IsPublished() {
		return false
	}

	if !indexableTypes[m.Type] {
		return false
	}

	return true
}

func (m *Meta) IsIndex() bool {
	if !m.IsPublished() {
		return false
	}

	return indexTypes[m.Type]
}

// ValidSectionOrDef lets easily check if a section is valid or return the default section.
func ValidSectionOrDef(section string) string {
	if section == "" {
		return DefSection
	}

	return section
}

// ValidTypeOrDef lets easily check if a content type is valid or return the default type.
func ValidTypeOrDef(contentType string) (defType string) {
	if !validTypes[contentType] {
		return ContentType.Article
	}

	return contentType
}

func (m Meta) CreatedAtPretty() string {
	if m.CreatedAt == "" {
		return "n/a"
	}

	createdAt, err := time.Parse(time.RFC3339, m.CreatedAt)
	if err != nil {
		return "n/a"
	}

	return createdAt.Format("January 2, 2006 15:04")
}

func (m Meta) PublishedAtPretty() string {
	if m.PublishedAt == "" {
		return "n/a"
	}

	publishedAt, err := time.Parse(time.RFC3339, m.PublishedAt)
	if err != nil {
		return "n/a"
	}

	return publishedAt.Format("January 2, 2006 15:04")
}

func (m Meta) UpdatedAtPretty() string {
	if m.UpdatedAt == "" {
		return "n/a"
	}

	updatedAt, err := time.Parse(time.RFC3339, m.UpdatedAt)
	if err != nil {
		return "n/a"
	}

	return updatedAt.Format("January 2, 2006 15:04")
}
