package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/inkcheck/readability"
)

var textExtensions = map[string]bool{
	".txt":  true,
	".md":   true,
	".rst":  true,
	".adoc": true,
	".tex":  true,
}

func main() {
	formula := flag.String("f", "", "readability formula to use")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: readability -f <formula> [files or directories...]\n")
		fmt.Fprintf(os.Stderr, "       command | readability -f <formula>\n\n")
		fmt.Fprintf(os.Stderr, "Available formulas:\n")
		for _, name := range readability.FormulaNames() {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *formula == "" {
		flag.Usage()
		os.Exit(1)
	}

	f := readability.Formula(*formula)
	if !f.Valid() {
		fmt.Fprintf(os.Stderr, "Error: unknown formula %q\n\nAvailable formulas:\n", *formula)
		for _, name := range readability.FormulaNames() {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
		os.Exit(1)
	}

	analyzer := &readability.Analyzer{}
	args := flag.Args()

	// No file args: read from stdin.
	if len(args) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		score, err := analyzer.Score(string(data), f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%.2f\n", score)
		return
	}

	// Collect files to process.
	var files []string
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if info.IsDir() {
			if err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && textExtensions[strings.ToLower(filepath.Ext(path))] {
					files = append(files, path)
				}
				return nil
			}); err != nil {
				fmt.Fprintf(os.Stderr, "Error walking %s: %v\n", arg, err)
				os.Exit(1)
			}
		} else {
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No files found.\n")
		os.Exit(1)
	}

	multiFile := len(files) > 1

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", path, err)
			continue
		}
		score, err := analyzer.Score(string(data), f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if multiFile {
			fmt.Printf("%s\t%.2f\n", path, score)
		} else {
			fmt.Printf("%.2f\n", score)
		}
	}
}
