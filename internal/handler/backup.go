package handler

import (
	"log"

	"github.com/adrianpk/gohermes/internal/hermes"
)

func Backup() error {
	config, err := hermes.LoadConfig()
	if err != nil {
		log.Fatalf("load config error: %v", err)
		return err
	}

	err = hermes.BackupToGitHub(config)
	if err != nil {
		log.Fatalf("backup error: %v", err)
		return err
	}

	log.Println("backup completed!")
	return nil
}
