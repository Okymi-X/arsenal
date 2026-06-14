package venv

import (
	"fmt"

	"github.com/Okymi-X/arsenal/internal/isolation"
)

// Install runs pip inside the virtualenv according to spec.
func (b *Backend) Install(spec isolation.InstallSpec) error {
	if !b.Exists() {
		return fmt.Errorf("virtualenv not provisioned at %s", b.dir)
	}
	args := buildPipArgs(spec)
	if len(args) == 0 {
		return fmt.Errorf("install spec produced no pip arguments")
	}
	if err := runCommand(b.pipExe(), args...); err != nil {
		return fmt.Errorf("pip install: %w", err)
	}
	return nil
}

// buildPipArgs translates a backend-agnostic spec into pip arguments.
func buildPipArgs(spec isolation.InstallSpec) []string {
	args := []string{"install", "--no-input"}
	args = append(args, spec.Extras...)
	if spec.GitURL != "" {
		args = append(args, gitTarget(spec)...)
	}
	args = append(args, spec.PipSpecs...)
	if len(args) == 2 {
		return nil
	}
	return args
}

func gitTarget(spec isolation.InstallSpec) []string {
	ref := spec.GitURL
	if spec.Commit != "" {
		ref = fmt.Sprintf("git+%s@%s", spec.GitURL, spec.Commit)
	} else {
		ref = "git+" + spec.GitURL
	}
	if spec.Editable {
		return []string{"-e", ref}
	}
	return []string{ref}
}
