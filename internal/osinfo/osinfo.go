// Package osinfo detects the host platform and architecture.
//
// It exposes a small, immutable snapshot of runtime facts so that other
// packages can make platform decisions without importing the runtime
// package directly or relying on global state.
package osinfo

import "runtime"

// Info is an immutable snapshot of host platform facts.
type Info struct {
	// OS is the GOOS value, e.g. "linux", "darwin".
	OS string
	// Arch is the GOARCH value, e.g. "amd64", "arm64".
	Arch string
}

// Detect returns the platform facts for the running host.
func Detect() Info {
	return Info{OS: runtime.GOOS, Arch: runtime.GOARCH}
}

// IsLinux reports whether the host runs Linux.
func (i Info) IsLinux() bool { return i.OS == "linux" }

// IsDarwin reports whether the host runs macOS.
func (i Info) IsDarwin() bool { return i.OS == "darwin" }

// String renders the platform as "os/arch".
func (i Info) String() string { return i.OS + "/" + i.Arch }
