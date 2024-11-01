package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	dirPath    string
	statsTypes string
	suffix     string
	overwrite  bool
	yes        bool
	workers    int
)

func init() {
	flag.StringVar(&dirPath, "dir", ".", "Directory to search for .pcap or .pcapng files")
	flag.StringVar(&statsTypes, "stats", "", "Statistics types (space-separated) or path to file containing types")
	flag.StringVar(&suffix, "suffix", ".stats-total.txt", "Suffix for resulting statistics files")
	flag.BoolVar(&overwrite, "overwrite", false, "Overwrite existing statistics files")
	flag.BoolVar(&yes, "yes", false, "Process files without asking for confirmation")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "Number of files to process in parallel")
}

func main() {
	flag.Parse()

	files, err := findPcapFiles(dirPath, overwrite, suffix)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No pcap or pcapng files need processing.")
		return
	}

	displayFiles(files)

	if !yes {
		if !confirm("Do you want to proceed? (y/n): ") {
			fmt.Println("Aborted.")
			return
		}
	}

	args, err := buildArgs(statsTypes)
	if err != nil {
		fmt.Printf("Error building arguments: %v\n", err)
		os.Exit(1)
	}

	processFiles(files, args)
}
