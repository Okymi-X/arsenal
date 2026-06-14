package doctor

import "strings"

// splitList splits a PATH-style list on the OS path list separator.
func splitList(list string) []string {
	return strings.Split(list, string(pathListSeparator))
}
