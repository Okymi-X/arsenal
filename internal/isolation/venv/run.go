package venv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Run executes a binary from the virtualenv's bin directory.
//
// The first element of args is the binary name; the remainder are passed
// through. Standard streams are connected to the calling process so the tool
// behaves as if run directly.
func (b *Backend) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command to run")
	}
	if !b.Exists() {
		return fmt.Errorf("virtualenv not provisioned at %s", b.dir)
	}
	bin := filepath.Join(b.binDir(), args[0])
	if _, err := os.Stat(bin); err != nil {
		return fmt.Errorf("binary %q not found in environment: %w", args[0], err)
	}
	cmd := exec.Command(bin, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "VIRTUAL_ENV="+b.dir, "PATH="+b.binDir()+string(os.PathListSeparator)+os.Getenv("PATH"))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", args[0], err)
	}
	return nil
}

// runCommand executes a command with stdio wired to the parent process.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
