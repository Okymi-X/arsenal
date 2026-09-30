package cli

import (
	"context"
	"fmt"

	"github.com/Okymi-X/arsenal/internal/installer"
	"github.com/Okymi-X/arsenal/internal/registry"
	"github.com/Okymi-X/arsenal/internal/resolver"
)

func resolveGitHubInstall(ctx context.Context, tool registry.Tool, ref string) (resolver.Resolved, error) {
	return resolveGitHubInstallExpected(ctx, tool, ref, "")
}

func resolveGitHubInstallExpected(ctx context.Context, tool registry.Tool, ref, expectedCommit string) (resolver.Resolved, error) {
	switch tool.InstallMethod {
	case installer.MethodPip, installer.MethodGitPip, installer.MethodGoBin, installer.MethodCargo:
	default:
		return resolver.Resolved{}, fmt.Errorf("direct GitHub install is not supported for %s tools", tool.InstallMethod)
	}
	repo, err := registry.ParseGitHubRepo(tool.Repo)
	if err != nil {
		return resolver.Resolved{}, err
	}
	client, err := registry.NewGitHubClient(repo, "registry.toml", githubToken())
	if err != nil {
		return resolver.Resolved{}, err
	}
	commit, err := client.ResolveCommit(ctx, ref)
	if err != nil {
		return resolver.Resolved{}, err
	}
	if err := verifySelectedGitHubCommit(ref, expectedCommit, commit); err != nil {
		return resolver.Resolved{}, err
	}
	if tool.InstallMethod == installer.MethodPip || tool.InstallMethod == installer.MethodGitPip {
		tool.InstallMethod = installer.MethodGitPip
	}
	return resolver.Resolved{
		Tool: tool,
		Version: registry.Version{
			Tag:            "github-" + commit[:12],
			Commit:         commit,
			Repo:           tool.Repo,
			Tested:         false,
			Notes:          "Explicit GitHub ref " + ref,
			InstallTargets: latestInstallTargets(tool),
		},
	}, nil
}

func latestInstallTargets(tool registry.Tool) map[string]string {
	for _, version := range tool.Versions {
		if len(version.InstallTargets) == 0 {
			continue
		}
		result := make(map[string]string, len(version.InstallTargets))
		for binary, target := range version.InstallTargets {
			result[binary] = target
		}
		return result
	}
	return nil
}
