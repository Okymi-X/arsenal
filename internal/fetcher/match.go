package fetcher

import "strings"

// normalizeName lowercases a file name and strips a common executable suffix so
// loose matching ignores case and the .exe/.sh/.bat extension.
func normalizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, ext := range []string{".exe", ".sh", ".bat"} {
		s = strings.TrimSuffix(s, ext)
	}
	return s
}

// matchName returns the candidate best matching want: an exact normalized match
// first, then a substring match. The bool is false when none matches.
func matchName(want string, candidates []string) (string, bool) {
	w := normalizeName(want)
	for _, c := range candidates {
		if normalizeName(c) == w {
			return c, true
		}
	}
	for _, c := range candidates {
		if strings.Contains(normalizeName(c), w) {
			return c, true
		}
	}
	return "", false
}
