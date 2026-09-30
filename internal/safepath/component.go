// Package safepath owns validation for names used as filesystem path components.
package safepath

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// ValidateComponent rejects empty, traversing, nested, or unsupported names.
func ValidateComponent(kind, value string) error {
	if value == "" || value == "." || value == ".." || filepath.Base(value) != value {
		return fmt.Errorf("invalid %s %q", kind, value)
	}
	for _, char := range value {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || strings.ContainsRune("._+-", char) {
			continue
		}
		return fmt.Errorf("invalid %s %q", kind, value)
	}
	return nil
}
