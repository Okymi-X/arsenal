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
// first, then the shortest substring match. Preferring the shortest substring
// makes a plain name win over a longer decorated sibling (for example the
// release zip over its "+debug" variant). The bool is false when none matches.
func matchName(want string, candidates []string) (string, bool) {
	w := normalizeName(want)
	for _, c := range candidates {
		if normalizeName(c) == w {
			return c, true
		}
	}
	best := ""
	for _, c := range candidates {
		if strings.Contains(normalizeName(c), w) {
			if best == "" || len(c) < len(best) {
				best = c
			}
		}
	}
	return best, best != ""
}
