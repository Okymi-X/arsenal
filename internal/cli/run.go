package cli

import (
	"fmt"
	"strings"
)

// cmdRun executes a binary from an installed tool's active version.
//
// Forms:
//
//	arsenal run <tool> -- <args...>            run the tool's primary binary
//	arsenal run <tool> <binary> -- <args...>   run a specific binary of the tool
//
// The binary selector is matched loosely, ignoring case, a ".py" suffix, and a
// "<tool>-" prefix, so "getTGT", "gettgt.py", and "impacket-getTGT" all match.
func (a *App) cmdRun(args []string) error {
	if len(args) == 0 {
		return usageError("run <tool> [binary] -- <args...>")
	}
	if strings.HasPrefix(args[0], "-") {
		return usageError("run <tool> [binary] -- <args...>")
	}
	name := a.canonicalName(args[0])

	m, err := a.store.Load()
	if err != nil {
		return err
	}
	active, ok := m.Active(name)
	if !ok {
		return fmt.Errorf("%s has no active version; run 'arsenal install %s'", name, name)
	}
	if len(active.Binaries) == 0 {
		return fmt.Errorf("%s has no runnable binary", name)
	}

	bin, forwarded := selectBinary(active.Name, active.Binaries, args[1:])
	backend := a.newBackend()
	if err := backend.Create(active.Name, active.Version); err != nil {
		return err
	}
	return backend.Run(append([]string{bin}, forwarded...))
}

// selectBinary resolves which binary to run and the arguments to forward.
//
// The arguments before a "--" separator may begin with a binary selector; if
// the first such token matches one of the tool's binaries it selects that
// binary, otherwise it is treated as a forwarded argument.
func selectBinary(tool string, binaries, rest []string) (bin string, forwarded []string) {
	pre, post := splitDashDash(rest)
	bin = binaries[0]
	if len(pre) > 0 {
		if match, ok := matchBinary(tool, binaries, pre[0]); ok {
			return match, append(append([]string{}, pre[1:]...), post...)
		}
	}
	return bin, append(append([]string{}, pre...), post...)
}

// splitDashDash splits args at the first "--" separator.
func splitDashDash(args []string) (before, after []string) {
	for i, a := range args {
		if a == "--" {
			return args[:i], args[i+1:]
		}
	}
	return args, nil
}

// matchBinary finds the binary whose normalized name equals the selector.
func matchBinary(tool string, binaries []string, selector string) (string, bool) {
	want := normalizeBinary(tool, selector)
	for _, b := range binaries {
		if normalizeBinary(tool, b) == want {
			return b, true
		}
	}
	return "", false
}

// normalizeBinary lowercases a name and strips a ".py" suffix and a "<tool>-"
// prefix so selectors match binaries loosely.
func normalizeBinary(tool, name string) string {
	n := strings.ToLower(strings.TrimSpace(name))
	n = strings.TrimSuffix(n, ".py")
	n = strings.TrimPrefix(n, strings.ToLower(tool)+"-")
	return n
}
