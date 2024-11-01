package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func findPcapFiles(dir string, overwrite bool, suffix string) ([]string, error) {
	var filesToProcess []string
	err := filepath.WalkDir(
		dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && (strings.HasSuffix(path, ".pcap") || strings.HasSuffix(path, ".pcapng")) {
				statsPath := strings.Replace(path, filepath.Ext(path), suffix, 1)
				if overwrite || !fileExists(statsPath) {
					filesToProcess = append(filesToProcess, path)
				}
			}
			return nil
		},
	)
	return filesToProcess, err
}
