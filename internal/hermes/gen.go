package hermes

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

func Parse(content []byte, path string) (Content, error) {
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

	md := parsed.Markdown

	renderer := NewTailwindRenderer(path)
	parsed.HTML = string(blackfriday.Run(md, blackfriday.WithRenderer(renderer)))
	return parsed, nil
}

// RenderNode renders a markdown node to HTML with Tailwind CSS classes.
// It handles various node types such as images, links, headings, paragraphs, lists, etc.
// At some point, the styles for each node type will be configurable.
func (r *TailwindRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	fmt.Printf("Rendering node: %v, entering: %v\n", node.Type, entering)
	switch node.Type {
	case blackfriday.Document:
		if entering {
			w.Write([]byte(`<style>
                .bullet-reset {
                    list-style: none;
                    padding: 0;
                    margin: 0;
                    position: relative;
                    padding-left: 20px;
                }
                .bullet-reset::before {
                    content: "\2022";
                    position: absolute;
                    left: 0;
                    top: 0;
                    color: black;
                }
            </style>`))
		}
	case blackfriday.Image:
		if entering {
			imagePath := string(node.LinkData.Destination)
			modifiedPath := modifyImagePath(imagePath, r.path)
			fmt.Printf("Modifying image path: %s -> %s\n", imagePath, modifiedPath)
			node.LinkData.Destination = []byte(modifiedPath)
			w.Write([]byte(`<img src="` + modifiedPath + `" alt="` + string(node.LinkData.Title) + `" class="max-w-full h-auto">`))
		}
	case blackfriday.Link:
		if entering {
			linkPath := string(node.LinkData.Destination)
			if strings.HasSuffix(linkPath, ".md") {
				modifiedPath := strings.TrimSuffix(linkPath, ".md") + ".html"
				fmt.Printf("Modifying link path: %s -> %s\n", linkPath, modifiedPath)
				node.LinkData.Destination = []byte(modifiedPath)
			}
			w.Write([]byte(`<a href="` + string(node.LinkData.Destination) + `" class="text-blue-500 underline">`))
		} else {
			w.Write([]byte(`</a>`))
		}
	case blackfriday.Heading:
		if entering {
			level := node.HeadingData.Level
			switch level {
			case 1:
				w.Write([]byte(`<h1 class="text-5xl font-bold mb-8 text-center text-gray-900">`))
			case 2:
				w.Write([]byte(`<h2 class="text-4xl font-bold mb-6 text-center text-gray-800">`))
			case 3:
				w.Write([]byte(`<h3 class="text-3xl font-bold mb-4 text-center text-gray-700">`))
			case 4:
				w.Write([]byte(`<h4 class="text-2xl font-bold mb-2 text-center text-gray-600">`))
			case 5:
				w.Write([]byte(`<h5 class="text-xl font-bold mb-1 text-center text-gray-500">`))
			case 6:
				w.Write([]byte(`<h6 class="text-lg font-bold mb-1 text-center text-gray-400">`))
			}
		} else {
			level := node.HeadingData.Level
			switch level {
			case 1:
				w.Write([]byte(`</h1>`))
			case 2:
				w.Write([]byte(`</h2>`))
			case 3:
				w.Write([]byte(`</h3>`))
			case 4:
				w.Write([]byte(`</h4>`))
			case 5:
				w.Write([]byte(`</h5>`))
			case 6:
				w.Write([]byte(`</h6>`))
			}
		}
	case blackfriday.Paragraph:
		if entering {
			w.Write([]byte(`<p class="text-lg text-gray-700 mb-4">`))
		} else {
			w.Write([]byte(`</p>`))
		}
	case blackfriday.List:
		if entering {
			w.Write([]byte(`<ul class="space-y-2">`))
		} else {
			w.Write([]byte(`</ul>`))
		}
	case blackfriday.Item:
		if entering {
			w.Write([]byte(`<li class="bullet-reset">`))
		} else {
			w.Write([]byte(`</li>`))
		}
	case blackfriday.BlockQuote:
		if entering {
			w.Write([]byte(`<blockquote class="border-l-4 border-gray-300 pl-4 italic text-gray-600">`))
		} else {
			w.Write([]byte(`</blockquote>`))
		}
	case blackfriday.HorizontalRule:
		if entering {
			w.Write([]byte(`<hr class="border-t-2 border-gray-300 my-4">`))
		}
	case blackfriday.Emph:
		if entering {
			w.Write([]byte(`<em>`))
		} else {
			w.Write([]byte(`</em>`))
		}
	case blackfriday.Strong:
		if entering {
			w.Write([]byte(`<strong>`))
		} else {
			w.Write([]byte(`</strong>`))
		}
	case blackfriday.Del:
		if entering {
			w.Write([]byte(`<del>`))
		} else {
			w.Write([]byte(`</del>`))
		}
	case blackfriday.Text:
		w.Write(node.Literal)
	case blackfriday.CodeBlock:
		if entering {
			w.Write([]byte(`<pre class="bg-gray-800 text-white p-4 rounded"><code class="language-go">`))
			w.Write(node.Literal)
			w.Write([]byte(`</code></pre>`))
		}
	case blackfriday.Code:
		if entering {
			w.Write([]byte(`<code class="bg-gray-200 p-1 rounded">`))
			w.Write(node.Literal)
			w.Write([]byte(`</code>`))
		}
	case blackfriday.Table:
		if entering {
			w.Write([]byte(`<table class="table-auto w-full">`))
		} else {
			w.Write([]byte(`</table>`))
		}
	case blackfriday.TableCell:
		if entering {
			w.Write([]byte(`<td class="border px-4 py-2">`))
		} else {
			w.Write([]byte(`</td>`))
		}
	case blackfriday.TableHead:
		if entering {
			w.Write([]byte(`<thead class="bg-gray-200">`))
		} else {
			w.Write([]byte(`</thead>`))
		}
	case blackfriday.TableBody:
		if entering {
			w.Write([]byte(`<tbody>`))
		} else {
			w.Write([]byte(`</tbody>`))
		}
	case blackfriday.TableRow:
		if entering {
			w.Write([]byte(`<tr>`))
		} else {
			w.Write([]byte(`</tr>`))
		}
	case blackfriday.Softbreak:
		w.Write([]byte("\n"))
	case blackfriday.Hardbreak:
		w.Write([]byte("<br>"))
	default:
		return r.Renderer.RenderNode(w, node, entering)
	}
	return blackfriday.GoToNext
}

type TailwindRenderer struct {
	blackfriday.Renderer
	path string
}

func NewTailwindRenderer(path string) *TailwindRenderer {
	return &TailwindRenderer{
		Renderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{}),
		path:     path,
	}
}

// Remove this one after demonstrating that logic in RenderNode works in the same way
func updatePaths(mdContent []byte, path string) []byte {
	log.Println("Starting updatePaths")

	node := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions)).Parse(mdContent)

	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering {
			switch node.Type {
			case blackfriday.Image:
				imagePath := string(node.LinkData.Destination)
				modifiedPath := modifyImagePath(imagePath, path)
				log.Printf("Modifying image path: %s -> %s", imagePath, modifiedPath)
				node.LinkData.Destination = []byte(modifiedPath)
			case blackfriday.Link:
				linkPath := string(node.LinkData.Destination)
				if strings.HasSuffix(linkPath, ".md") {
					modifiedPath := strings.TrimSuffix(linkPath, ".md") + ".html"
					log.Printf("Modifying link path: %s -> %s", linkPath, modifiedPath)
					node.LinkData.Destination = []byte(modifiedPath)
				}
			}
		}
		return blackfriday.GoToNext
	})

	var buf bytes.Buffer
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})
	renderer.RenderHeader(&buf, node)
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		renderer.RenderNode(&buf, node, entering)
		return blackfriday.GoToNext
	})
	renderer.RenderFooter(&buf, node)

	log.Println("Finished updatePaths")
	return buf.Bytes()
}

// Remove this one after demonstrating that logic in RenderNode works in the same way
func updateImgPaths(mdContent []byte, path string) []byte {
	node := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions)).Parse(mdContent)

	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Image && entering {
			imagePath := string(node.LinkData.Destination)
			modifiedPath := modifyImagePath(imagePath, path)
			node.LinkData.Destination = []byte(modifiedPath)
		}
		return blackfriday.GoToNext
	})

	var buf bytes.Buffer
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})
	renderer.RenderHeader(&buf, node)
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		renderer.RenderNode(&buf, node, entering)
		return blackfriday.GoToNext
	})
	renderer.RenderFooter(&buf, node)

	return buf.Bytes()
}

func modifyImagePath(imagePath, mdPath string) string {
	imageDir := strings.TrimSuffix(mdPath, filepath.Ext(mdPath))
	relativeImageDir := strings.TrimPrefix(imageDir, ContentDir+"/")

	if strings.HasPrefix(relativeImageDir, "root/") {
		relativeImageDir = strings.TrimPrefix(relativeImageDir, "root/")
		parts := strings.Split(relativeImageDir, string(os.PathSeparator))

		if len(parts) == 1 {
			relativeImageDir = filepath.Join(ImgDir, parts[0])
		} else if len(parts) > 1 {
			switch parts[0] {
			case ContentType.Blog, ContentType.Series:
				relativeImageDir = filepath.Join(ImgDir, parts[0], parts[1])
			case ContentType.Article, ContentType.Page:
				relativeImageDir = filepath.Join(ImgDir, parts[1])
			default:
				relativeImageDir = filepath.Join(ImgDir, strings.Join(parts, string(os.PathSeparator)))
			}
		}
	} else {
		parts := strings.Split(relativeImageDir, string(os.PathSeparator))

		if len(parts) > 2 && (parts[1] == ContentType.Article || parts[1] == ContentType.Page) {
			relativeImageDir = filepath.Join(ImgDir, parts[0], parts[2])
		} else if len(parts) > 2 && (parts[1] == ContentType.Blog || parts[1] == ContentType.Series) {
			relativeImageDir = filepath.Join(ImgDir, parts[0], parts[1], parts[2])
		} else {
			relativeImageDir = filepath.Join(ImgDir, strings.Join(parts, string(os.PathSeparator)))
		}
	}

	return filepath.Join(relativeImageDir, filepath.Base(imagePath))
}

// Remove this one after demonstrating that logic in RenderNode works in the same way
func updateDocPaths(mdContent []byte) []byte {
	log.Println("Starting updateDocPaths")

	log.Printf("Original Markdown content:\n%s", string(mdContent))

	node := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions)).Parse(mdContent)

	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Link && entering {
			linkPath := string(node.LinkData.Destination)
			log.Printf("Found link: %s", linkPath)
			if strings.HasSuffix(linkPath, ".md") {
				modifiedPath := strings.TrimSuffix(linkPath, ".md") + ".html"
				log.Printf("Modifying link: %s -> %s", linkPath, modifiedPath)
				node.LinkData.Destination = []byte(modifiedPath)
			} else {
				log.Printf("Link does not end with .md: %s", linkPath)
			}
		} else {
			log.Printf("Node type: %v, entering: %v", node.Type, entering)
		}
		return blackfriday.GoToNext
	})

	var buf bytes.Buffer
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})
	renderer.RenderHeader(&buf, node)
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		renderer.RenderNode(&buf, node, entering)
		return blackfriday.GoToNext
	})
	renderer.RenderFooter(&buf, node)

	log.Println("Finished updateDocPaths")
	return buf.Bytes()
}
