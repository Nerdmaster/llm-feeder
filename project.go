package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
)

type fsEntry struct {
	name     string
	path     string
	relative string
	isDir    bool
	entries  []*fsEntry
}

func (e *fsEntry) sortEntries() {
	sort.Slice(e.entries, func(i, j int) bool {
		return e.entries[i].name < e.entries[j].name
	})
}

type project struct {
	ignores []string
	root    *fsEntry
}

func newProject(base string) (*project, error) {
	var err error
	base, err = filepath.Abs(base)
	if err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}
	var i os.FileInfo
	i, err = os.Stat(base)
	if err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}
	if !i.IsDir() {
		return nil, fmt.Errorf("invalid project path: not a directory")
	}

	return &project{root: &fsEntry{name: "<root>", path: base, isDir: true}}, nil
}

// shouldIgnore checks if an entry should be ignored
func (p *project) shouldIgnore(pth string) bool {
	for _, pattern := range p.ignores {
		var matched, _ = filepath.Match(pattern, pth)
		if matched {
			return true
		}
	}
	return false
}

// scan looks for all dir entries inside e's path and adds them recursively to
// its list, using the given ignore function to potentially skip some files
func (e *fsEntry) scan(ignore func(string) bool) error {
	var entries, err = os.ReadDir(e.path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		var name = entry.Name()
		var fullPath = filepath.Join(e.path, name)
		var relPath = filepath.Join(e.relative, name)
		if ignore(relPath) || ignore(name) {
			slog.Debug("ignoring entry per ignore list", "path", fullPath)
			continue
		}

		var newEntry = &fsEntry{name: name, path: fullPath, relative: relPath}

		slog.Debug("adding item", "entry", newEntry.path, "parent", e.path)
		e.entries = append(e.entries, newEntry)

		if entry.IsDir() {
			newEntry.isDir = true
			newEntry.scan(ignore)
		}
	}

	e.sortEntries()
	return nil
}

// scanAll looks for all files that don't match one of the ignore patterns, and
// puts it into the project entries
func (p *project) scanAll() error {
	return p.root.scan(p.shouldIgnore)
}
