package hermes

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	outputDir = "output"
	repo    = "https://github.com/adrianpk/adrianpk.github.io"
	pubBranch  = "gh-pages"
	mainBranch = "main"
)

type Options struct {
	RepoURL   string
	PubBranch string
	MainBranch string
}

func PublishToGitHubPages(cfg Config) error {
	repoURL := cfg.RepoURL()
	pubBranch := cfg.Repo.Pub
	mainBranch := cfg.Repo.Main

	err := os.Chdir(outputDir)
	if err != nil {
		return err
	}

	_, err = os.Stat(".git")
	if os.IsNotExist(err) {
		_, _, err = runCommand("git", "init")
		if err != nil {
			return err
		}
		_, _, err = runCommand("git", "remote", "add", "origin", repoURL)
		if err != nil {
			return err
		}
	} else {
		stdout, stderr, err := runCommand("git", "status", "--porcelain")
		if err != nil {
			return fmt.Errorf("git status error: %s: stdout: %s, stderr: %s", err, stdout, stderr)
		}
		if stdout != "" {
			log.Println("Uncommitted changes detected, committing them.")
		}
	}

	_, _, err = runCommand("git", "add", ".")
	if err != nil {
		return err
	}

	stdout, stderr, err := runCommand("git", "commit", "-m", "Deploy to GitHub Pages")
	if err != nil && !strings.Contains(stderr, "nothing to commit") {
		return fmt.Errorf("git commit error: %s: stdout: %s, stderr: %s", err, stdout, stderr)
	}

	stdout, stderr, err = runCommand("git", "pull", "origin", pubBranch, "--rebase")
	if err != nil && strings.Contains(stderr, fmt.Sprintf("couldn't find remote ref %s", pubBranch)) {
		log.Printf("Remote branch %s does not exist, skipping pull step.", pubBranch)
	} else if err != nil {
		return fmt.Errorf("git pull error: %s: stdout: %s, stderr: %s", err, stdout, stderr)
	}

	_, _, err = runCommand("git", "push", "origin", mainBranch+":"+pubBranch)
	if err != nil {
		return err
	}

	err = os.Chdir("..")
	if err != nil {
		return err
	}

	return nil
}

func runCommand(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String(), stderr.String(), fmt.Errorf("%s: %s", err, stderr.String())
	}
	return out.String(), stderr.String(), nil
}
