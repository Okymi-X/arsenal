package venv

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Okymi-X/arsenal/internal/safepath"
)

// Run executes a binary from the virtualenv's bin directory.
//
// The first element of args is the binary name; the remainder are passed
// through. Standard streams are connected to the calling process so the tool
// behaves as if run directly.
func (b *Backend) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command to run")
	}
	if err := safepath.ValidateComponent("binary", args[0]); err != nil {
		return err
	}
	if !b.Exists() {
		return fmt.Errorf("virtualenv not provisioned at %s", b.dir)
	}
	bin := filepath.Join(b.binDir(), args[0])
	if _, err := os.Stat(bin); err != nil {
		return fmt.Errorf("binary %q not found in environment: %w", args[0], err)
	}
	cmd := exec.CommandContext(ctx, bin, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = b.environment()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", args[0], err)
	}
	return nil
}

func runCommandEnv(ctx context.Context, env []string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}
