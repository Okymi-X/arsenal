package cli

import (
	"fmt"

	"github.com/Okymi-X/arsenal/internal/registry"
	"github.com/Okymi-X/arsenal/internal/resolver"
)

func (a *App) resolveInstallOptions(reg *registry.Registry, options installOptions) (resolver.Resolved, error) {
	var err error
	if options.selectTool {
		options.spec, err = a.selectRegistryTool(reg, options.query)
		if err != nil {
			return resolver.Resolved{}, err
		}
	}
	if options.selectVersion {
		options.spec, err = a.selectCuratedVersion(reg, options.spec)
		if err != nil {
			return resolver.Resolved{}, err
		}
	}
	resolved, err := resolveSpec(reg, options.spec)
	if err != nil {
		return resolver.Resolved{}, err
	}
	if options.selectGitHub {
		selected, selectErr := a.selectToolGitHubRef(resolved.Tool)
		if selectErr != nil {
			return resolver.Resolved{}, selectErr
		}
		return resolveGitHubInstallExpected(a.ctx, resolved.Tool, selected.Name, selected.Commit)
	}
	if options.githubRef != "" {
		return resolveGitHubInstall(a.ctx, resolved.Tool, options.githubRef)
	}
	return resolved, nil
}

func (a *App) selectRegistryTool(reg *registry.Registry, query string) (string, error) {
	matches := reg.Search(query)
	if len(matches) == 0 {
		return "", fmt.Errorf("no tools match %q", query)
	}
	options := make([]selectOption, 0, len(matches))
	for _, tool := range matches {
		options = append(options, selectOption{
			label: fmt.Sprintf("%-16s [%s] %s", tool.Name, tool.Category, tool.Description),
			value: tool.Name,
		})
	}
	return selectOne(a.in, a.out, "Select a tool:", options)
}

func (a *App) selectCuratedVersion(reg *registry.Registry, name string) (string, error) {
	tool, err := reg.MustFindTool(name)
	if err != nil {
		return "", err
	}
	options := make([]selectOption, 0, len(tool.Versions))
	for _, version := range tool.Versions {
		status := "untested"
		if version.Tested {
			status = "tested"
		}
		options = append(options, selectOption{
			label: fmt.Sprintf("%-24s %s", version.Tag, status),
			value: version.Tag,
		})
	}
	selected, err := selectOne(a.in, a.out, "Select a curated version for "+tool.Name+":", options)
	if err != nil {
		return "", err
	}
	return tool.Name + "@" + selected, nil
}

func (a *App) selectToolGitHubRef(tool registry.Tool) (registry.GitHubRef, error) {
	repo, err := registry.ParseGitHubRepo(tool.Repo)
	if err != nil {
		return registry.GitHubRef{}, err
	}
	client, err := registry.NewGitHubClient(repo, "registry.toml", githubToken())
	if err != nil {
		return registry.GitHubRef{}, err
	}
	refs, err := client.ListTags(a.ctx, 100)
	if err != nil {
		return registry.GitHubRef{}, err
	}
	options := make([]selectOption, 0, len(refs))
	byName := make(map[string]registry.GitHubRef, len(refs))
	for _, ref := range refs {
		byName[ref.Name] = ref
		status := "upstream"
		if version, ok := curatedGitHubVersion(tool, ref.Name); ok {
			status = "curated"
			if version.Tested {
				status = "curated, tested"
			}
		}
		options = append(options, selectOption{
			label: fmt.Sprintf("%-24s %s  %s", ref.Name, ref.Commit[:12], status),
			value: ref.Name,
		})
	}
	selected, err := selectOne(a.in, a.out, "Select an upstream GitHub tag for "+tool.Name+":", options)
	if err != nil {
		return registry.GitHubRef{}, err
	}
	return byName[selected], nil
}

func curatedGitHubVersion(tool registry.Tool, ref string) (registry.Version, bool) {
	for _, version := range tool.Versions {
		if version.Tag == ref || version.Commit == ref || "v"+version.Tag == ref {
			return version, true
		}
	}
	return registry.Version{}, false
}
