package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func DeleteTempFiles(commentary []string) {
	// Define the directory where your files are located

	for id, _ := range commentary {
		//directory := "/path/to/your/directory"
		// Call a function to delete specific file types
		deleteFilesByExtension(fmt.Sprintf("audio-%d.mp3", id), ".mp3") // Delete audio files with ".mp3" extension
		deleteFilesByExtension(fmt.Sprintf("image-%d.png", id), ".png") // Delete image files with ".jpg" extension
		deleteFilesByExtension(fmt.Sprintf("chunk-%d.mp4", id), ".mp4") // Delete video files with ".mp4" extension
	}

}

func deleteFilesByExtension(directory, extension string) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error walking directory:", err)
			return err
		}

		// Check if the file has the specified extension
		if !info.IsDir() && filepath.Ext(path) == extension {
			err := os.Remove(path)
			if err != nil {
				fmt.Println("Error deleting file:", err)
				return err
			}

			fmt.Printf("Deleted: %s\n", path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}
