package fm

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func Preprocess(content []byte) (Content, error) {
	var parsed Content

	if bytes.HasPrefix(content, []byte("---\n")) {
		parts := bytes.SplitN(content, []byte("---\n"), 3)
		if len(parts) == 3 {
			err := yaml.Unmarshal(parts[1], &parsed.Metadata)
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

	return parsed, nil
}
