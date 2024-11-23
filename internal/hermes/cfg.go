package hermes

import (
	"fmt"
	"html/template"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	cfgFile   = "hermes.yml"
	githubURL = "https://github.com/%s/%s"
)

type Config struct {
	Name     string     `yaml:"name"`
	Git      GitConfig  `yaml:"git"`
	PubRepo  RepoConfig `yaml:"pubRepo"`
	BakRepo  RepoConfig `yaml:"bakRepo"`
	Menu     []string   `yaml:"menu"`
	Sections []Section  `yaml:"sections"`
}

type GitConfig struct {
	User string `yaml:"user"`
}

type RepoConfig struct {
	Name   string `yaml:"name"`
	Main   string `yaml:"main"`
	Pub    string `yaml:"pub"`
	Update string `yaml:"update"`
}

type Section struct {
	Name         string   `yaml:"name"`
	ContentTypes []string `yaml:"content_types"`
}

const cfgYAML = `name: {{.Name}}
git:
  user: {{.Git.User}}
pubRepo:
  main: {{.PubRepo.Main}}
  pub: {{.PubRepo.Pub}}
bakRepo:
  name: {{.BakRepo.Name}}
  main: {{.BakRepo.Main}}
  update: {{.BakRepo.Update}}
menu:
{{- range .Menu }}
  - {{ . }}
{{- end }}
sections:
{{- range .Sections }}
  - name: {{ .Name }}
    content_types:
    {{- range .ContentTypes }}
      - {{ . }}
    {{- end }}
{{- end }}`

func NewCfgFile(name string, user string) error {
	file, err := os.Create("hermes.yml")
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New("cfg").Parse(cfgYAML)
	if err != nil {
		return err
	}

	data := Config{
		Name: name,
		Git: GitConfig{
			User: user,
		},
		PubRepo: RepoConfig{
			Main: "main",
			Pub:  "gh-pages",
		},
		BakRepo: RepoConfig{
			Name:   "site.hermes.3",
			Main:   "main",
			Update: "new/wip",
		},
		Menu: []string{"about-us", "contact"},
		Sections: []Section{
			{
				Name:         "root",
				ContentTypes: []string{"articles", "blog", "series"},
			},
			{
				Name:         "section",
				ContentTypes: []string{"articles", "blog", "series"},
			},
		},
	}

	return tmpl.Execute(file, data)
}

func LoadConfig() (Config, error) {
	var config Config
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

func (cfg *Config) PubRepoURL() string {
	return fmt.Sprintf(githubURL, cfg.Git.User, cfg.PubRepo.Name)
}

func (cfg *Config) BakRepoURL() string {
	return fmt.Sprintf(githubURL, cfg.Git.User, cfg.BakRepo.Name)
}
