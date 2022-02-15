package config

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

const (
	workspaceFilename     = "WORKSPACE"
	aspectpluginsFilename = ".aspectplugins"
)

// Finder is the interface that wraps the simple Find method that performs the
// finding of the plugins file in the user system.
type Finder interface {
	Find() (string, error)
}

type finder struct {
	osGetwd func() (string, error)
	osStat  func(string) (fs.FileInfo, error)
}

// NewFinder instantiates a default internal implementation of the Finder
// interface.
func NewFinder() Finder {
	return &finder{
		osGetwd: os.Getwd,
		osStat:  os.Stat,
	}
}

// Find finds the .aspectplugins file under a Bazel workspace. If the returned
// path is empty and no error was produced, the file doesn't exist.
func (f *finder) Find() (string, error) {
	cwd, err := f.osGetwd()
	if err != nil {
		return "", fmt.Errorf("failed to find .aspectplugins: %w", err)
	}
	for {
		if cwd == "/" {
			break
		}
		workspacePath := path.Join(cwd, workspaceFilename)
		if _, err := f.osStat(workspacePath); err != nil {
			if !os.IsNotExist(err) {
				return "", fmt.Errorf("failed to find .aspectplugins: %w", err)
			}
			cwd = filepath.Dir(cwd)
			continue
		}
		aspectpluginsPath := path.Join(cwd, aspectpluginsFilename)
		if _, err := f.osStat(aspectpluginsPath); err != nil {
			if !os.IsNotExist(err) {
				return "", fmt.Errorf("failed to find .aspectplugins: %w", err)
			}
			break
		}
		return aspectpluginsPath, nil
	}
	return "", nil
}
