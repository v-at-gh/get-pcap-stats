package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

func processFiles(files []string, args []string) {
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	totalFiles := len(files)

	for i, file := range files {
		wg.Add(1)
		sem <- struct{}{}

		go func(f string, idx int) {
			defer wg.Done()
			defer func() { <-sem }()
			processFile(f, idx, totalFiles, args)
		}(file, i)
	}

	wg.Wait()
	close(sem)
}

func processFile(filePath string, index int, totalFiles int, args []string) {
	statsPath := strings.Replace(filePath, filepath.Ext(filePath), suffix, 1)
	fmt.Printf("Processing %d/%d: %s...\n", index+1, totalFiles, filePath)

	if err := extractStats(filePath, statsPath, args); err != nil {
		fmt.Printf("Error processing %s: %v\n", filePath, err)
	}
}
