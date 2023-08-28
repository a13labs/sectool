package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func listKeys(parentPath string) ([]string, error) {
	var subfoldersWithIDFiles []string

	entries, err := os.ReadDir(parentPath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subfolderPath := filepath.Join(parentPath, entry.Name())

			files, err := os.ReadDir(subfolderPath)
			if err != nil {
				continue
			}

			for _, file := range files {
				if !file.IsDir() && strings.HasPrefix(file.Name(), "id_") {
					subfoldersWithIDFiles = append(subfoldersWithIDFiles, entry.Name())
					break // No need to continue checking files in this subfolder
				}
			}
		}
	}

	return subfoldersWithIDFiles, nil
}
