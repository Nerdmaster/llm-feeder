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

// createProject populates the project data with metadata about its base path
func createProject(base string, ignores []string) (*project, error) {
	var p = &project{ignores: ignores, root: &fsEntry{name: "<root>", path: base, isDir: true}}
	slog.Debug(fmt.Sprintf("project created: %#v", p))
	var err = p.scan(p.root)
	return p, err
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

func (p *project) scan(e *fsEntry) error {
	var entries, err = os.ReadDir(e.path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		var name = entry.Name()
		var fullPath = filepath.Join(e.path, name)
		var relPath = filepath.Join(e.relative, name)
		if p.shouldIgnore(relPath) || p.shouldIgnore(name) {
			slog.Debug("ignoring entry per ignore list", "path", fullPath)
			continue
		}

		var newEntry = &fsEntry{name: name, path: fullPath, relative: relPath}

		slog.Debug("adding item", "entry", newEntry.path, "parent", e.path)
		e.entries = append(e.entries, newEntry)

		if entry.IsDir() {
			newEntry.isDir = true
			p.scan(newEntry)
		}
	}

	e.sortEntries()
	return nil
}
