package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteOptions controls file writing behavior.
type WriteOptions struct {
	Force  bool
	Backup bool
	Yes    bool
}

// Writer handles file creation with safety policies.
type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

// Write writes content to filename in the current directory.
func (w *Writer) Write(filename, content string, opts WriteOptions) error {
	// Validate filename (path traversal protection)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	target := filepath.Join(cwd, filename)
	clean := filepath.Clean(target)
	if !strings.HasPrefix(clean, filepath.Clean(cwd)) {
		return fmt.Errorf("invalid filename: path traversal detected")
	}

	// Check if file exists
	_, err = os.Stat(clean)
	fileExists := err == nil

	if fileExists {
		if opts.Force && opts.Backup {
			if err := w.createBackup(clean); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
		} else if opts.Backup {
			if err := w.createBackup(clean); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
		} else if opts.Force {
			// proceed to overwrite
		} else if opts.Yes {
			// proceed silently
		} else {
			// Ask for confirmation
			confirmed, err := confirm(fmt.Sprintf("File %s already exists. Overwrite?", filename))
			if err != nil {
				return err
			}
			if !confirmed {
				return fmt.Errorf("operation cancelled by user")
			}
		}
	}

	if err := os.WriteFile(clean, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	return nil
}

func (w *Writer) createBackup(path string) error {
	bakPath := path + ".bak"
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := os.WriteFile(bakPath, data, 0644); err != nil {
		return err
	}
	fmt.Printf("✔ Backup created: %s\n", filepath.Base(bakPath))
	return nil
}

func confirm(question string) (bool, error) {
	fmt.Printf("%s [y/N]: ", question)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes", nil
}
