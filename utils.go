package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		return text == "y" || text == "Y"
	}
	return false
}

func displayFiles(files []string) {
	fmt.Println("The following files will be processed:")
	for i, f := range files {
		fmt.Printf("%d: %s\n", i+1, f)
	}
	fmt.Printf("Total %d files.\n", len(files))
}
