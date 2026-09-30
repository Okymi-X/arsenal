// Package pipspec validates exact Python package requirements used by the registry.
package pipspec

import (
	"fmt"
	"strings"
)

// Pin is an exact package name and version pair.
type Pin struct {
	Name    string
	Version string
}

// NormalizedName returns the canonical comparison form defined for Python
// project names by replacing separator runs with a single hyphen.
func (p Pin) NormalizedName() string {
	var out strings.Builder
	separator := false
	for _, char := range strings.ToLower(p.Name) {
		if strings.ContainsRune("._-", char) {
			if !separator {
				out.WriteByte('-')
			}
			separator = true
			continue
		}
		separator = false
		out.WriteRune(char)
	}
	return out.String()
}

// ParseExact accepts the restricted name==version form used for reproducible
// registry dependencies. Options, URLs, markers, and ranges are rejected.
func ParseExact(value string) (Pin, error) {
	if value == "" || strings.TrimSpace(value) != value {
		return Pin{}, fmt.Errorf("invalid exact pip requirement %q", value)
	}
	name, version, found := strings.Cut(value, "==")
	if !found || !validName(name) || !validVersion(version) {
		return Pin{}, fmt.Errorf("invalid exact pip requirement %q", value)
	}
	return Pin{Name: name, Version: version}, nil
}

func validName(value string) bool {
	if value == "" || !isAlphaNumeric(rune(value[0])) || !isAlphaNumeric(rune(value[len(value)-1])) {
		return false
	}
	for _, char := range value {
		if isAlphaNumeric(char) || strings.ContainsRune("._-", char) {
			continue
		}
		return false
	}
	return true
}

func validVersion(value string) bool {
	if value == "" || !isAlphaNumeric(rune(value[0])) {
		return false
	}
	for _, char := range value {
		if isAlphaNumeric(char) || strings.ContainsRune(".!+_-", char) {
			continue
		}
		return false
	}
	return true
}

func isAlphaNumeric(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9'
}
