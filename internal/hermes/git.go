package hermes

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	outputDir  = "output"
	repo       = "https://github.com/adrianpk/adrianpk.github.io"
	pubBranch  = "gh-pages"
	mainBranch = "main"
)


func PublishToGitHubPages(cfg Config) error {
	repoURL := fmt.Sprintf("https://github.com/%s/%s", cfg.Git.User, cfg.PubRepo.Name)
	pubBranch := cfg.PubRepo.Pub
	mainBranch := cfg.PubRepo.Main

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

func BackupToGitHub(cfg Config) error {
	repoURL := fmt.Sprintf("https://github.com/%s/%s", cfg.Git.User, cfg.BakRepo.Name)
	mainBranch := cfg.BakRepo.Main
	updateBranch := cfg.BakRepo.Update

	err := os.Chdir(".")
	if err != nil {
		return err
	}

	_, err = os.Stat(".git")
	if os.IsNotExist(err) {
		log.Println("Initializing new git repository")
		_, _, err = runCommand("git", "init")
		if err != nil {
			return err
		}
		_, _, err = runCommand("git", "remote", "add", "origin", repoURL)
		if err != nil {
			return err
		}
		_, _, err = runCommand("git", "checkout", "-b", mainBranch)
		if err != nil {
			return err
		}
		_, _, err = runCommand("git", "add", "--all", ":!output")
		if err != nil {
			return err
		}

		commitMessage := fmt.Sprintf("Initial commit at %s", time.Now().Format(time.RFC3339))
		stdout, stderr, err := runCommand("git", "commit", "-m", commitMessage)
		if err != nil && !strings.Contains(stderr, "nothing to commit") {
			return fmt.Errorf("git commit error: %s: stdout: %s, stderr: %s", err, stdout, stderr)
		}

		log.Println("Pushing initial commit to main branch")
		stdout, stderr, err = runCommand("git", "push", "origin", mainBranch)
		if err != nil {
			log.Printf("Error pushing to main branch: %s\nstdout: %s\nstderr: %s", err, stdout, stderr)
			return err
		}

		log.Println("Creating and switching to update branch")
		_, _, err = runCommand("git", "checkout", "-b", updateBranch)
		if err != nil {
			return err
		}

		log.Println("Pushing initial commit to update branch")
		stdout, stderr, err = runCommand("git", "push", "origin", updateBranch)
		if err != nil {
			log.Printf("Error pushing to update branch: %s\nstdout: %s\nstderr: %s", err, stdout, stderr)
			return err
		}
	} else {
		log.Println("Switching to update branch")
		_, _, err = runCommand("git", "checkout", updateBranch)
		if err != nil {
			return err
		}

		log.Println("Adding changes")
		_, _, err = runCommand("git", "add", "--all", ":!output")
		if err != nil {
			return err
		}

		commitMessage := fmt.Sprintf("Update at %s", time.Now().Format(time.RFC3339))
		stdout, stderr, err := runCommand("git", "commit", "-m", commitMessage)
		if err != nil && !strings.Contains(stderr, "nothing to commit") {
			return fmt.Errorf("git commit error: %s: stdout: %s, stderr: %s", err, stdout, stderr)
		}

		log.Println("Pushing changes to update branch")
		stdout, stderr, err = runCommand("git", "push", "origin", updateBranch)
		if err != nil {
			log.Printf("Error pushing to update branch: %s\nstdout: %s\nstderr: %s", err, stdout, stderr)
			return err
		}
	}

	err = os.Chdir("..")
	if err != nil {
		return err
	}

	log.Println("Successfully backed up to GitHub")
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
