package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// OrganizeFiles moves files into folders based on their extensions
func OrganizeFiles(sourceDir string) error {
	// Map file extensions to folder names
	folderMap := map[string]string{
		".jpg": "Images",
		".png": "Images",
		".gif": "Images",
		".mp4": "Videos",
		".mkv": "Videos",
		".txt": "Documents",
		".pdf": "Documents",
	}

	// Walk through files in the source directory
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get file extension
		ext := filepath.Ext(info.Name())
		if folder, ok := folderMap[ext]; ok {
			// Create destination folder if it doesn't exist
			destDir := filepath.Join(sourceDir, folder)
			if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create folder: %v", err)
			}

			// Move file
			destPath := filepath.Join(destDir, info.Name())
			if err := os.Rename(path, destPath); err != nil {
				return fmt.Errorf("failed to move file: %v", err)
			}

			fmt.Printf("Moved %s -> %s\n", path, destPath)
		}
		return nil
	})
}
