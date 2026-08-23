package engine

import (
	"os"
	"path/filepath"
)

func DataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".prediction-engine"
	}
	return filepath.Join(home, ".prediction-engine")
}

func Stat(p string) (os.FileInfo, error) {
	return os.Stat(p)
}
