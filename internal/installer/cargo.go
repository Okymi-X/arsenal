package installer

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Okymi-X/arsenal/internal/registry"
)

var cargoVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`)

// CargoMethod installs exact crates into isolated roots.
type CargoMethod struct {
	toolsRoot string
	runner    CommandRunner
}

// NewCargoMethod returns a Cargo installer with explicit filesystem and process dependencies.
func NewCargoMethod(toolsRoot string, runner CommandRunner) *CargoMethod {
	return &CargoMethod{toolsRoot: toolsRoot, runner: runner}
}

// Supports reports whether the tool declares the cargo install method.
func (m *CargoMethod) Supports(tool registry.Tool) bool {
	return tool.InstallMethod == MethodCargo
}

// Install builds exact crates in staging and promotes verified binaries.
func (m *CargoMethod) Install(ctx context.Context, tool registry.Tool, version registry.Version) (Result, error) {
	directGitHub := strings.HasPrefix(version.Tag, "github-") && commitSHA.MatchString(version.Commit)
	if !directGitHub && !cargoVersion.MatchString(version.Tag) {
		return Result{}, fmt.Errorf("cargo version %q is not an exact semantic version", version.Tag)
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
	for _, target := range targets {
		command, err := cargoInstallCommand(stage, tool, version, target, directGitHub)
		if err != nil {
			return Result{}, err
		}
		if err := m.installTarget(ctx, stage, target, command); err != nil {
			return Result{}, err
		}
	}
	if err := removeInstallTree(filepath.Join(stage, "cache")); err != nil {
		return Result{}, fmt.Errorf("clean Cargo install cache: %w", err)
	}
	if err := verifyBinaries(stage, tool.AllBinaries()); err != nil {
		return Result{}, err
	}
	if err := promoteStaging(stage, final); err != nil {
		return Result{}, err
	}
	return Result{Path: final, Backend: MethodCargo}, nil
}

func (m *CargoMethod) installTarget(ctx context.Context, stage string, target installTarget, command Command) error {
	if err := validateCargoTarget(target.Source); err != nil {
		return err
	}
	if err := m.runner.Run(ctx, command); err != nil {
		return fmt.Errorf("install Cargo crate %s: %w", target.Source, err)
	}
	return renameInstalledBinary(stage, target.Source, target.Binary)
}

func cargoInstallCommand(stage string, tool registry.Tool, version registry.Version, target installTarget, direct bool) (Command, error) {
	args := []string{"install", "--root", stage}
	if direct {
		repo := tool.RepoFor(version)
		if _, err := registry.ParseGitHubRepo(repo); err != nil {
			return Command{}, err
		}
		args = append(args, "--git", repo, "--rev", version.Commit)
	} else {
		args = append(args, "--version", version.Tag)
	}
	args = append(args, "--locked", target.Source)
	return Command{
		Name: "cargo",
		Args: args,
		Env: map[string]string{
			"CARGO_HOME":       filepath.Join(stage, "cache", "cargo"),
			"CARGO_TARGET_DIR": filepath.Join(stage, "cache", "target"),
		},
	}, nil
}

func validateCargoTarget(target string) error {
	if target == "" || strings.Trim(target, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-") != "" {
		return fmt.Errorf("invalid Cargo crate %q", target)
	}
	return nil
}

var _ InstallMethod = (*CargoMethod)(nil)
