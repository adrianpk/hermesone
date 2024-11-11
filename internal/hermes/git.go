package hermes

import (
	"os"
	"os/exec"
)

const (
	outputDir = "output"
	repoURL   = "https://github.com/adrianpk/sample-site.git"
	branch    = "gh-pages"
)

func PublishToGitHubPages() error {
	if err := os.Chdir(outputDir); err != nil {
		return err
	}

	if err := exec.Command("git", "init").Run(); err != nil {
		return err
	}

	if err := exec.Command("git", "remote", "add", "origin", repoURL).Run(); err != nil {
		return err
	}

	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return err
	}

	if err := exec.Command("git", "commit", "-m", "Deploy to GitHub Pages").Run(); err != nil {
		return err
	}

	if err := exec.Command("git", "push", "-f", "origin", "master:"+branch).Run(); err != nil {
		return err
	}

	if err := os.RemoveAll(".git"); err != nil {
		return err
	}

	if err := os.Chdir(".."); err != nil {
		return err
	}

	return nil
}
