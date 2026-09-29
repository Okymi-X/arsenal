package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Okymi-X/arsenal/internal/registry"
)

var (
	goMajorSuffix = regexp.MustCompile(`^v[2-9][0-9]*$`)
	goVersionRef  = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`)
	commitSHA     = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

// GoBinMethod installs pinned Go command packages into isolated roots.
type GoBinMethod struct {
	toolsRoot string
	runner    CommandRunner
}

// NewGoBinMethod returns a Go installer with explicit filesystem and process dependencies.
func NewGoBinMethod(toolsRoot string, runner CommandRunner) *GoBinMethod {
	return &GoBinMethod{toolsRoot: toolsRoot, runner: runner}
}

// Supports reports whether the tool declares the gobin install method.
func (m *GoBinMethod) Supports(tool registry.Tool) bool {
	return tool.InstallMethod == MethodGoBin
}

// Install builds each pinned command in staging and promotes verified binaries.
func (m *GoBinMethod) Install(ctx context.Context, tool registry.Tool, version registry.Version) (Result, error) {
	if !goVersionRef.MatchString(version.Commit) && !commitSHA.MatchString(version.Commit) {
		return Result{}, fmt.Errorf("gobin version %q has no exact install ref", version.Tag)
	}
	targets, err := orderedTargets(tool, version)
	if err != nil {
		return Result{}, err
	}
	stage, final, err := prepareStaging(m.toolsRoot, tool.Name, version.Tag)
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = removeInstallTree(stage) }()
	if err := os.Mkdir(filepath.Join(stage, "bin"), 0o755); err != nil {
		return Result{}, fmt.Errorf("create Go binary directory: %w", err)
	}
	for _, target := range targets {
		if err := m.installTarget(ctx, stage, version.Commit, target); err != nil {
			return Result{}, err
		}
	}
	if err := removeInstallTree(filepath.Join(stage, "cache")); err != nil {
		return Result{}, fmt.Errorf("clean Go install cache: %w", err)
	}
	if err := verifyBinaries(stage, tool.AllBinaries()); err != nil {
		return Result{}, err
	}
	if err := promoteStaging(stage, final); err != nil {
		return Result{}, err
	}
	return Result{Path: final, Backend: MethodGoBin}, nil
}

func (m *GoBinMethod) installTarget(ctx context.Context, stage, version string, target installTarget) error {
	if err := validateGoTarget(target.Source); err != nil {
		return err
	}
	command := Command{
		Name: "go",
		Args: []string{"install", target.Source + "@" + version},
		Env: map[string]string{
			"GOBIN":      filepath.Join(stage, "bin"),
			"GOCACHE":    filepath.Join(stage, "cache", "build"),
			"GOMODCACHE": filepath.Join(stage, "cache", "mod"),
			"GOPATH":     filepath.Join(stage, "cache", "gopath"),
		},
	}
	if err := m.runner.Run(ctx, command); err != nil {
		return fmt.Errorf("install Go package %s: %w", target.Source, err)
	}
	return renameInstalledBinary(stage, goBinaryName(target.Source), target.Binary)
}

func validateGoTarget(target string) error {
	if target == "" || strings.ContainsAny(target, "@\\ \t\r\n") || strings.Contains(target, "://") {
		return fmt.Errorf("invalid Go install target %q", target)
	}
	parts := strings.Split(target, "/")
	if len(parts) < 3 || !strings.Contains(parts[0], ".") {
		return fmt.Errorf("invalid Go install target %q", target)
	}
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return fmt.Errorf("invalid Go install target %q", target)
		}
	}
	return nil
}

func goBinaryName(target string) string {
	name := filepath.Base(target)
	if goMajorSuffix.MatchString(name) {
		name = filepath.Base(filepath.Dir(target))
	}
	return name
}

var _ InstallMethod = (*GoBinMethod)(nil)
