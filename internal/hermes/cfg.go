package hermes

import (
	"html/template"
	"os"
)

const cfgYAML = `name: {{.Name}}`

func NewCfgFile(name string) error {
	file, err := os.Create("hermes.yml")
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New("cfg").Parse(cfgYAML)
	if err != nil {
		return err
	}

	data := struct {
		Name string
	}{
		Name: name,
	}

	return tmpl.Execute(file, data)
}
