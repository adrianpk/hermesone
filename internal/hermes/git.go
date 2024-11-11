package hermes

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const (
	outputDir = "output"
	repoURL   = "https://github.com/adrianpk/sample-site.git"
	branch    = "gh-pages"
)

func PublishToGitHubPages() error {
	err := os.Chdir(outputDir)
	if err != nil {
		return err
	}

	_, err = os.Stat(".git")
	if os.IsNotExist(err) {
		_, err = runCommand("git", "init")
		if err != nil {
			return err
		}
		_, err = runCommand("git", "remote", "add", "origin", repoURL)
		if err != nil {
			return err
		}
	}

	_, err = runCommand("git", "add", ".")
	if err != nil {
		return err
	}

	out, err := runCommand("git", "commit", "-m", "Deploy to GitHub Pages")
	if err != nil {
		return fmt.Errorf("git commit error: %s: %s", err, out)
	}

	_, err = runCommand("git", "push", "-f", "origin", "main:"+branch)
	if err != nil {
		return err
	}

	err = os.Chdir("..")
	if err != nil {
		return err
	}

	return nil
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("%s: %s", err, stderr.String())
	}
	return out.String(), nil
}
