package util

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadProjectEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return godotenv.Load(filepath.Join(dir, "..", ".env"))
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return errors.New("go.mod not found")
		}
		dir = parent
	}
}
