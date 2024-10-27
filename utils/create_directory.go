package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateAuthorDirectory(authorName string, baseDir string) error {
	dir := filepath.Join(baseDir, authorName)
	subDirs := []string{"thumb_mini", "small", "regular", "original"}

	for _, subDir := range subDirs {
		fullPath := filepath.Join(dir, subDir)
		err := os.MkdirAll(fullPath, 0755) // MkdirAll will create all parent directories as needed
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return err
		}
	}
	fmt.Println("\n"+"\033[1;34m"+"Successfully created directory structure for: ", authorName, "\033[0m")
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
