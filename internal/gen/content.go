package gen

type Content struct {
	Meta     Meta
	Markdown []byte
	HTML     string
}

type Meta struct {
	Title           string
	Description     string
	Summary         string
	Date            string
	PublishDate     string
	LastModified    string
	Type            string
	Section         string
	Slug            string
	Image           string
	SocialImage     string
	Layout          string
	CanonicalURL    string
	Locale          string
	Robots          string
	Excerpt         string
	Permalink       string
	Draft           bool
	TableOfContents bool
	Share           bool
	Featured        bool
	Comments        bool
	Authors         []string
	Categories      []string
	Tags            []string
	Keywords        []string
	Sitemap         Sitemap
}

type Sitemap struct {
	Priority   float64
	ChangeFreq string
}
