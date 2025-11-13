package util

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sirupsen/logrus"
)

// ProcessTemplate processes a template string with the given arguments
func ProcessTemplate(tmplStr string, args map[string]interface{}) string {
	logrus.Debugf("Processing template: %s", tmplStr)
	tmpl, err := template.New("stateName").Parse(tmplStr)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, args)
	if err != nil {
		logrus.Errorf("Failed to execute template: %v", err)
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return buf.String()
}

// EnsureFileExists creates a directory if it doesn't exist and returns a file handle for appending
func EnsureFileExists(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logrus.Debugf("Creating directory: %s", dir)
		// Create directory if it doesn't exist
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logrus.Errorf("Failed to create directory %s: %v", dir, err)
			return nil, err
		}
	}

	// Create or open the file for appending
	logrus.Debugf("Opening file: %s", path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("Failed to open file %s: %v", path, err)
		return nil, err
	}
	return file, nil
}