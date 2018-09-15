package reader

import (
	"log"
	"os"
	"path/filepath"
)

var extWhitelist = []string{".vue", ".yml", ".yaml"}

func Walk(dir string) []string {
	log.Printf("Collecting %v...", dir)

	var filesToBeProcessed []string

	err := filepath.Walk(dir, func(location string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("%v cannot be accessed", location)
		}

		if processRequired(location) {
			filesToBeProcessed = append(filesToBeProcessed, location)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walkthrough the location %q: %v", dir, err)
	}

	return filesToBeProcessed
}

func processRequired(path string) bool {
	ext := filepath.Ext(path)

	for _, v := range extWhitelist {
		if ext == v {
			return true
		}
	}

	return false
}
