package hermes

import "path/filepath"

const (
	ContentDir   = "content"
	OutputDir    = "output"
	LayoutDir    = "assets/layout"
	DefLayoutDir = "default"
	ImgDir       = "img"
	DefSection   = "root"
	IndexMdFile  = "index.md"
	IndexFile    = "index.html"
	DefLayout    = "default.html"
	DefHeaderImg = "header.png"
	NoJekyllFile = ".nojekyll"
)

var (
	DefLayoutPath       = filepath.Join(LayoutDir, DefLayoutDir)
	DefLayoutFile       = filepath.Join(DefLayoutPath, DefLayout)
	BakLayoutPathFormat = filepath.Join(DefLayoutPath, "default-%s.html.bak")
)

var (
	separator = []byte("---\n")
)

type contentType struct {
	Page    string
	Article string
	Blog    string
	Series  string
}

var ContentType = contentType{
	Page:    "page",
	Article: "article",
	Blog:    "blog",
	Series:  "series",
}
