package internal

import (
	"log"
	"os"
)

func InitFile() {
	// Define the filename
	filename := "data.csv"

	// Check if the file exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, delete it
		err := os.Remove(filename)
		if err != nil {
			log.Println("Error deleting file:", err)
			return
		}
		log.Println("File deleted:", filename)
	}

	// Create a new csv file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error creating file:", err)
		return
	}
	defer file.Close()

}
