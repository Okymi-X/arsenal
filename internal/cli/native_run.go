package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Okymi-X/arsenal/internal/installer"
	"github.com/Okymi-X/arsenal/internal/store"
)

func (a *App) runNative(tool store.InstalledTool, binary string, args []string) error {
	expected, err := installer.InstallationPath(a.paths.Tools, tool.Name, tool.Version)
	if err != nil {
		return err
	}
	if filepath.Clean(tool.Path) != expected {
		return fmt.Errorf("refusing to run %s from untrusted path %q", tool.Name, tool.Path)
	}
	target := filepath.Join(expected, "bin", binary)
	info, err := os.Lstat(target)
	if err != nil {
		return fmt.Errorf("inspect %s binary: %w", tool.Name, err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
		return fmt.Errorf("refusing to run non-executable binary %q", target)
	}
	cmd := exec.CommandContext(a.ctx, target, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = nativeEnvironment(filepath.Join(expected, "bin"))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", binary, err)
	}
	return nil
}

func nativeEnvironment(binDir string) []string {
	environment := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "PATH=") {
			environment = append(environment, entry)
		}
	}
	return append(environment, "PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}
