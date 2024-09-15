package hermes

type Content struct {
	Meta     Meta
	Markdown []byte
	HTML     string
}

type Meta struct {
	Title           string   `yaml:"title"`
	Description     string   `yaml:"description"`
	Summary         string   `yaml:"summary"`
	Date            string   `yaml:"date"`
	PublishDate     string   `yaml:"publish-date"`
	LastModified    string   `yaml:"last-modified"`
	Type            string   `yaml:"type"`
	Section         string   `yaml:"section"`
	Slug            string   `yaml:"slug"`
	Image           string   `yaml:"image"`
	SocialImage     string   `yaml:"social-image"`
	Layout          string   `yaml:"layout"`
	CanonicalURL    string   `yaml:"canonical-url"`
	Locale          string   `yaml:"locale"`
	Robots          string   `yaml:"robots"`
	Excerpt         string   `yaml:"excerpt"`
	Permalink       string   `yaml:"permalink"`
	Draft           bool     `yaml:"draft"`
	TableOfContents bool     `yaml:"table-of-contents"`
	Share           bool     `yaml:"share"`
	Featured        bool     `yaml:"featured"`
	Comments        bool     `yaml:"comments"`
	Authors         []string `yaml:"author"`
	Categories      []string `yaml:"categories"`
	Tags            []string `yaml:"tags"`
	Keywords        []string `yaml:"keywords"`
	Sitemap         Sitemap  `yaml:"sitemap"`
}

type Sitemap struct {
	Priority   float64 `yaml:"priority"`
	ChangeFreq string  `yaml:"changefreq"`
}
