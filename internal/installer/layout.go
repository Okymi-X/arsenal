package installer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Okymi-X/arsenal/internal/safepath"
)

// InstallationPath returns the only valid root for a tool/version pair.
func InstallationPath(toolsRoot, tool, version string) (string, error) {
	if err := safepath.ValidateComponent("tool", tool); err != nil {
		return "", err
	}
	if err := safepath.ValidateComponent("version", version); err != nil {
		return "", err
	}
	if toolsRoot == "" {
		return "", fmt.Errorf("tools root is empty")
	}
	return filepath.Join(toolsRoot, tool, version), nil
}

// RemoveInstallation deletes only the validated root recorded for a tool version.
func RemoveInstallation(toolsRoot, tool, version, recordedPath string) error {
	expected, err := InstallationPath(toolsRoot, tool, version)
	if err != nil {
		return err
	}
	if filepath.Clean(recordedPath) != expected {
		return fmt.Errorf("refusing to remove untrusted install path %q", recordedPath)
	}
	if err := removeInstallTree(expected); err != nil {
		return fmt.Errorf("remove installation %s: %w", expected, err)
	}
	return nil
}

func prepareStaging(toolsRoot, tool, version string) (string, string, error) {
	final, err := InstallationPath(toolsRoot, tool, version)
	if err != nil {
		return "", "", err
	}
	parent := filepath.Dir(final)
	if err := ensureSafeParent(toolsRoot, parent); err != nil {
		return "", "", err
	}
	stage, err := os.MkdirTemp(parent, "."+version+".install-")
	if err != nil {
		return "", "", fmt.Errorf("create install staging directory: %w", err)
	}
	return stage, final, nil
}

func ensureSafeParent(toolsRoot, parent string) error {
	if err := os.MkdirAll(toolsRoot, 0o755); err != nil {
		return fmt.Errorf("create tools root: %w", err)
	}
	if info, err := os.Lstat(parent); err == nil {
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return fmt.Errorf("tool directory is not a regular directory: %s", parent)
		}
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect tool directory: %w", err)
	}
	if err := os.Mkdir(parent, 0o755); err != nil {
		return fmt.Errorf("create tool directory: %w", err)
	}
	return nil
}

func promoteStaging(stage, final string) error {
	backup := stage + ".previous"
	hadPrevious := false
	if _, err := os.Lstat(final); err == nil {
		if err := os.Rename(final, backup); err != nil {
			return fmt.Errorf("preserve previous installation: %w", err)
		}
		hadPrevious = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect installation destination: %w", err)
	}
	if err := os.Rename(stage, final); err != nil {
		if hadPrevious {
			_ = os.Rename(backup, final)
		}
		return fmt.Errorf("activate staged installation: %w", err)
	}
	if hadPrevious {
		if err := removeInstallTree(backup); err != nil {
			return fmt.Errorf("remove replaced installation: %w", err)
		}
	}
	return nil
}

func removeInstallTree(root string) error {
	if _, err := os.Lstat(root); os.IsNotExist(err) {
		return nil
	}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return os.Chmod(path, 0o700)
		}
		return nil
	}); err != nil {
		return err
	}
	return os.RemoveAll(root)
}

func verifyBinaries(root string, binaries []string) error {
	if len(binaries) == 0 {
		return fmt.Errorf("installation exposes no binaries")
	}
	for _, binary := range binaries {
		if err := safepath.ValidateComponent("binary", binary); err != nil {
			return err
		}
		info, err := os.Lstat(filepath.Join(root, "bin", binary))
		if err != nil {
			return fmt.Errorf("installed binary %q: %w", binary, err)
		}
		if !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
			return fmt.Errorf("installed binary %q is not an executable regular file", binary)
		}
	}
	return nil
}
