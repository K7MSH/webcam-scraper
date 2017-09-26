package main

import (
	"os"
	"strings"
)

func ensureDir(path string) error {
	logger.Tracef("Ensuring dir '%s' exists", path)
	if !strings.Contains(path, string(os.PathSeparator)) {
		logger.Tracef("Path doesn't contain a %c, bailing", os.PathSeparator)
		return nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Tracef("Path doesn't exist. Attempting to create")
		err = os.MkdirAll(path, os.ModePerm)
		return err
	}
	logger.Tracef("Path exists already")
	return nil
}
