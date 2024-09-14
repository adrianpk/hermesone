package fm

type Content struct {
	Metadata Metadata
	Markdown []byte
}

type Metadata struct {
	Title           string   `yaml:"title"`
	Description     string   `yaml:"description"`
	Authors         []string `yaml:"author"`
	Date            string   `yaml:"date"`
	PublishDate     string   `yaml:"publish-date"`
	LastModified    string   `yaml:"last-modified"`
	Draft           bool     `yaml:"draft"`
	Type            string   `yaml:"type"`
	Section         string   `yaml:"section"`
	Categories      []string `yaml:"categories"`
	Tags            []string `yaml:"tags"`
	Keywords        []string `yaml:"keywords"`
	Slug            string   `yaml:"slug"`
	Summary         string   `yaml:"summary"`
	Image           string   `yaml:"image"`
	SocialImage     string   `yaml:"social-image"`
	Layout          string   `yaml:"layout"`
	CanonicalURL    string   `yaml:"canonical-url"`
	TableOfContents bool     `yaml:"table-of-contents"`
	Share           bool     `yaml:"share"`
	Featured        bool     `yaml:"featured"`
	Locale          string   `yaml:"locale"`
	Sitemap         struct {
		Priority   float64 `yaml:"priority"`
		ChangeFreq string  `yaml:"changefreq"`
	} `yaml:"sitemap"`
	Robots    string `yaml:"robots"`
	Excerpt   string `yaml:"excerpt"`
	Permalink string `yaml:"permalink"`
	Comments  bool   `yaml:"comments"`
}
