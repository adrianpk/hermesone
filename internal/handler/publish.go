package handler

import (
	"log"

	"github.com/adrianpk/gohermes/internal/hermes"
)

func Publish() error {
	err := hermes.PublishToGitHubPages()
	if err != nil {
		log.Fatalf("publish error: %v", err)
		return err
	}

	log.Println("published to GitHub Pages")
	return nil
}

