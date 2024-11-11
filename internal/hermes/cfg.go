package hermes

import (
	"fmt"
	"html/template"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	cfgFile = "hermes.yml"
)

type Config struct {
	Name string     `yaml:"name"`
	Repo RepoConfig `yaml:"repo"`
}

type RepoConfig struct {
	User string `yaml:"user"`
	Name string `yaml:"name"`
	Main string `yaml:"main"`
	Pub  string `yaml:"pub"`
}

// cfgYAML is the template for the hermes.yml configuration file.
// name is the name of the project.
// repo is the repository information.
// user is the GitHub username.
// name is the repository name.
// main is the main branch.
// pub is the branch to publish to.
const cfgYAML = `name: {{.Name}}
repo:
  user: {{.Repo.User}}
  name: {{.Repo.Name}}
  main: {{.Repo.Main}}
  pub: {{.Repo.Pub}}`

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
		Repo: RepoConfig{
			User: user,
			Name: user + ".github.io",
			Main: "main",
			Pub:  "gh-pages",
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

func (cfg *Config) RepoURL() string {
	return fmt.Sprintf("https://github.com/%s/%s", cfg.Repo.User, cfg.Repo.Name)
}
