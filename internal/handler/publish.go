package handler

import (
	"log"

	"github.com/adrianpk/gohermes/internal/hermes"
)

func Publish() error {
	config, err := hermes.LoadConfig()
	if err != nil {
		return err
	}

	err = hermes.PublishToGitHubPages(config)
	if err != nil {
		return err
	}

	log.Println("published!")

	return nil
}
