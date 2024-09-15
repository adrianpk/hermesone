package hermes

import (
	"bytes"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

func Parse(content []byte) (Content, error) {
	var parsed Content

	if bytes.HasPrefix(content, []byte("---\n")) {
		parts := bytes.SplitN(content, []byte("---\n"), 3)
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

	parsed.HTML = string(blackfriday.Run(parsed.Markdown))

	return parsed, nil
}
