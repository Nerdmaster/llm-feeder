package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log/slog"
	"os"
	"strings"
)

type stringSlice struct {
	values []string
}

func (s *stringSlice) String() string {
	return strings.Join(s.values, ", ")
}

func (s *stringSlice) Set(value string) error {
	s.values = append(s.values, value)
	return nil
}

func main() {
	// Define command-line flags
	var ignores stringSlice
	var base string
	var dirsOnly, verbose bool
	flag.Var(&ignores, "I", "Specify a shell pattern to ignore (can be used multiple times)")
	flag.Var(&ignores, "ignore", "Specify a shell pattern to ignore (can be used multiple times)")
	flag.BoolVar(&dirsOnly, "q", false, "Quiet (only list directories, not file contents)")
	flag.StringVar(&base, "d", ".", "Base directory to use (defaults to current working dir)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-i <pattern to ignore>] [-d <base directory>] [-q]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	var logOpts = slog.HandlerOptions{Level: slog.LevelInfo}
	if verbose {
		logOpts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &logOpts)))

	var p, err = createProject(base, ignores.values)
	if err != nil {
		slog.Error("createProject failure", "directory", base, "error", err)
		os.Exit(1)
	}

	// Print out tree view first
	fmt.Printf("Contents of %s:\n", base)
	printTree(p.root, -1)

	if dirsOnly {
		os.Exit(0)
	}

	printFiles(p.root)
}

func printTree(entry *fsEntry, indent int) {
	if indent != -1 {
		fmt.Println(strings.Repeat("   ", indent) + "-- " + entry.name)
	}
	if entry.isDir {
		for _, child := range entry.entries {
			printTree(child, indent+1)
		}
	}
}

func printFiles(entry *fsEntry) {
	if entry.isDir {
		for _, child := range entry.entries {
			printFiles(child)
		}
		return
	}

	var out = fileContents(entry.path)
	if out != "" {
		fmt.Printf("----------- BEGIN Contents of %q:\n", entry.relative)
		fmt.Println(out)
		fmt.Printf("----------- END Contents of %q\n\n", entry.relative)
	}
}

func fileContents(fname string) string {
	var data, err = ioutil.ReadFile(fname)
	if err != nil {
		slog.Error("Unable to read file", "file", fname, "error", err)
		os.Exit(1)
	}
	var test = data
	if len(data) > 10240 {
		test = data[:10240]
	}
	var i = bytes.IndexByte(test, 0)
	if i != -1 {
		return "<binary data skipped>"
	}

	return string(data)
}
