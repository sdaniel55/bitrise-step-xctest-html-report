package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/bitrise-io/go-utils/log"
)

func copy(sourcePath, outputDir string, errors *[]error) string {
	source, err := os.Open(sourcePath)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("Failed to open file, error: %s", err))
	}

	defer func() {
		if cerr := source.Close(); cerr != nil {
			*errors = append(*errors, fmt.Errorf("Failed to close file, error: %s", cerr))
		}
	}()

	destinationPath := path.Join(outputDir, path.Base(sourcePath))
	destination, err := os.Create(destinationPath)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("Failed to open file, error: %s", err))
	}

	defer func() {
		if cerr := destination.Close(); cerr != nil {
			*errors = append(*errors, fmt.Errorf("Failed to close file, error: %s", cerr))
		}
	}()

	if _, err := io.Copy(destination, source); err != nil {
		*errors = append(*errors, fmt.Errorf("Failed to copy file, error: %s", err))
	}

	return destinationPath
}

func installedInPath(name string) bool {
	cmd := exec.Command("which", name)
	outBytes, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(outBytes)) != ""
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	log.Warnf("For more details you can enable the debug logs by turning on the verbose step input.")
	os.Exit(1)
}
