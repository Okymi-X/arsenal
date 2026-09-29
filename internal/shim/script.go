package shim

import (
	"fmt"
	"strings"
)

// script renders the shell shim that execs the target binary.
//
// The shim is intentionally minimal: it forwards all arguments and replaces
// the process so signals and exit codes propagate unchanged.
func script(binary, targetBin string) string {
	return fmt.Sprintf(`#!/bin/sh
# arsenal shim for %s
# This file is generated. Do not edit; run 'arsenal switch' instead.
exec %s "$@"
`, binary, shellQuote(targetBin))
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}
