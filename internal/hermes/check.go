package hermes

import (
	"errors"
	"os"
)

var ErrNoHermes = errors.New("not an hermes project")

// CheckHermes checks if the current directory is an hermes project.
func CheckHermes() error {
	if _, err := os.Stat("hermes.yml"); os.IsNotExist(err) {
		return ErrNoHermes
	}
	return nil
}
