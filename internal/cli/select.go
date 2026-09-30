package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const maxSelectionInput = 1024

type selectOption struct {
	label string
	value string
}

// selectOne prints a bounded numbered menu and accepts only a numeric choice.
// Values are never parsed from displayed labels, keeping untrusted metadata out
// of the command-selection boundary.
func selectOne(in *bufio.Reader, out io.Writer, heading string, options []selectOption) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no choices are available")
	}
	if len(options) == 1 {
		return options[0].value, nil
	}
	if in == nil || out == nil {
		return "", fmt.Errorf("interactive selection requires standard input and output")
	}
	fmt.Fprintln(out, safeDisplay(heading, 160))
	for i, option := range options {
		fmt.Fprintf(out, "  %d) %s\n", i+1, safeDisplay(option.label, 240))
	}
	fmt.Fprintf(out, "Select [1-%d]: ", len(options))

	line, err := in.ReadSlice('\n')
	if errors.Is(err, bufio.ErrBufferFull) {
		return "", fmt.Errorf("selection input exceeds %d bytes", maxSelectionInput)
	}
	if err != nil && len(line) == 0 {
		return "", fmt.Errorf("selection canceled")
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read selection: %w", err)
	}
	if len(line) > maxSelectionInput {
		return "", fmt.Errorf("selection input exceeds %d bytes", maxSelectionInput)
	}
	choice, err := strconv.Atoi(strings.TrimSpace(string(line)))
	if err != nil || choice < 1 || choice > len(options) {
		return "", fmt.Errorf("invalid selection (want 1-%d)", len(options))
	}
	return options[choice-1].value, nil
}

func safeDisplay(value string, limit int) string {
	var result strings.Builder
	for _, character := range value {
		if result.Len() >= limit {
			result.WriteString("...")
			break
		}
		if unicode.IsControl(character) {
			result.WriteByte('?')
			continue
		}
		result.WriteRune(character)
	}
	return result.String()
}
