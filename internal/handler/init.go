package handler

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	osFileSep = string(os.PathSeparator)
)

// InitDirs initializes the directory structure and copies assets from the embedded filesystem.
func InitDirs(dirs []string, assetsFS embed.FS) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	err := copyDir(assetsFS, "assets", "assets")
	if err != nil {
		return err
	}

	return nil
}

// copyDir copies a directory from the embedded filesystem to the target location.
func copyDir(assetsFS embed.FS, srcDir, destDir string) error {
	return fs.WalkDir(assetsFS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		content, err := assetsFS.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, content, 0644)
	})
}
