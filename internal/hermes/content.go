package hermes

type Content struct {
	Meta     Meta
	Markdown []byte
	HTML     string
}

type Meta struct {
	Title           string   `yaml:"title"`         // Title of the content
	Description     string   `yaml:"description"`   // Description of the content
	Summary         string   `yaml:"summary"`       // Summary of the content
	Date            string   `yaml:"date"`          // Date of the content
	PublishDate     string   `yaml:"publish-date"`  // Publish date of the content
	LastModified    string   `yaml:"last-modified"` // Last modified date of the content
	Type            string   `yaml:"type"`          // Type of the content
	Section         string   `yaml:"section"`       // Section of the content
	Slug            string   `yaml:"slug"`          // Slug of the content
	Image           string   `yaml:"image"`         // Image of the content
	SocialImage     string   `yaml:"social-image"`  // Social image of the content
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
