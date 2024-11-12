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
	Name    string     `yaml:"name"`
	Git     GitConfig  `yaml:"git"`
	PubRepo RepoConfig `yaml:"pubRepo"`
	BakRepo RepoConfig `yaml:"bakRepo"`
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

// cfgYAML is the template for the hermes.yml configuration file.
// name is the name of the project.
// repo is the repository information.
// user is the GitHub username.
// name is the repository name.
// main is the main branch.
// pub is the branch to publish to.
const cfgYAML = `name: {{.Name}}
git:
  user: {{.Git.User}}
pubRepo:
  main: {{.PubRepo.Main}}
  pub: {{.PubRepo.Pub}}
bakRepo:
  name: {{.BakRepo.Name}}
  main: {{.BakRepo.Main}}
  update: {{.BakRepo.Update}}`

// NewCfgFile creates a new configuration file with default values for repo name, main branch, and pub branch.
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
			Name:   "something",
			Main:   "main",
			Update: "update",
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
