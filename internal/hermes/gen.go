package hermes

import (
	"bytes"
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

	md := updateImgPaths(parsed.Markdown, path)

	parsed.HTML = string(blackfriday.Run(md))

	return parsed, nil
}

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
			relativeImageDir = filepath.Join("img", parts[0])
		} else if len(parts) > 1 {
			switch parts[0] {
			case "blog", "series":
				relativeImageDir = filepath.Join("img", parts[0], parts[1])
			case "articles", "pages":
				relativeImageDir = filepath.Join("img", parts[1])
			default:
				relativeImageDir = filepath.Join("img", strings.Join(parts, string(os.PathSeparator)))
			}
		}
	} else {
		parts := strings.Split(relativeImageDir, string(os.PathSeparator))

		if len(parts) > 2 && (parts[1] == "articles" || parts[1] == "pages") {
			relativeImageDir = filepath.Join("img", parts[0], parts[2])
		} else if len(parts) > 2 && (parts[1] == "blog" || parts[1] == "series") {
			relativeImageDir = filepath.Join("img", parts[0], parts[1], parts[2])
		} else {
			relativeImageDir = filepath.Join("img", strings.Join(parts, string(os.PathSeparator)))
		}
	}

	return filepath.Join(relativeImageDir, filepath.Base(imagePath))
}
