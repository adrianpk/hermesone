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
	outputDir = "output"
)

func PublishToGitHubPages(cfg Config) error {
	repoURL := cfg.PubRepoURL()
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
			log.Println("uncommitted changes detected, committing them.")
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
		log.Printf("remote branch %s does not exist, skipping pull step.", pubBranch)
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
	repoURL := cfg.BakRepoURL()
	mainBranch := cfg.BakRepo.Main
	updateBranch := cfg.BakRepo.Update

	err := os.Chdir(".")
	if err != nil {
		return err
	}

	_, err = os.Stat(".git")
	if os.IsNotExist(err) {
		log.Println("initializing new git repository")
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

		log.Println("pushing initial commit to main branch")
		stdout, stderr, err = runCommand("git", "push", "origin", mainBranch)
		if err != nil {
			log.Printf("error pushing to main branch: %s\nstdout: %s\nstderr: %s", err, stdout, stderr)
			return err
		}

		log.Println("creating and switching to update branch")
		_, _, err = runCommand("git", "checkout", "-b", updateBranch)
		if err != nil {
			return err
		}

		log.Println("pushing initial commit to update branch")
		stdout, stderr, err = runCommand("git", "push", "origin", updateBranch)
		if err != nil {
			log.Printf("error pushing to update branch: %s\nstdout: %s\nstderr: %s", err, stdout, stderr)
			return err
		}
	} else {
		log.Println("switching to update branch")
		_, _, err = runCommand("git", "checkout", updateBranch)
		if err != nil {
			return err
		}

		log.Println("adding changes")
		_, _, err = runCommand("git", "add", "--all", ":!output")
		if err != nil {
			return err
		}

		commitMessage := fmt.Sprintf("update at %s", time.Now().Format(time.RFC3339))
		stdout, stderr, err := runCommand("git", "commit", "-m", commitMessage)
		if err != nil && !strings.Contains(stderr, "nothing to commit") {
			return fmt.Errorf("git commit error: %s: stdout: %s, stderr: %s", err, stdout, stderr)
		}

		log.Println("pushing changes to update branch")
		stdout, stderr, err = runCommand("git", "push", "origin", updateBranch)
		if err != nil {
			log.Printf("error pushing to update branch: %s\nstdout: %s\nstderr: %s", err, stdout, stderr)
			return err
		}
	}

	err = os.Chdir("..")
	if err != nil {
		return err
	}

	log.Println("backup complete!")
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
