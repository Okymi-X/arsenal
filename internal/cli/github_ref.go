package cli

import (
	"fmt"
	"strings"
)

func verifySelectedGitHubCommit(ref, expected, resolved string) error {
	if expected == "" || strings.EqualFold(expected, resolved) {
		return nil
	}
	return fmt.Errorf(
		"GitHub ref %q moved from %s to %s; review the available tags and select again",
		ref,
		shortCommit(expected),
		shortCommit(resolved),
	)
}

func shortCommit(commit string) string {
	if len(commit) <= 12 {
		return commit
	}
	return commit[:12]
}
