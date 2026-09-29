package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/Okymi-X/arsenal/internal/registry"
)

type syncOptions struct {
	listRefs  bool
	reference string
	repo      string
}

// cmdSync refreshes the local registry from its configured URL or a selected
// GitHub ref resolved to an immutable commit.
func (a *App) cmdSync(args []string) error {
	options, err := parseSyncArgs(args)
	if err != nil {
		return err
	}
	repo := a.cfg.RegistryRepo
	if options.repo != "" {
		repo = options.repo
	}
	if options.listRefs {
		return a.listGitHubRefs(repo)
	}
	source, sourceURL, pinned, err := a.syncSource(options, repo)
	if err != nil {
		return err
	}
	if err := a.paths.EnsureDirs(); err != nil {
		return err
	}
	a.log.Printf("-> syncing registry from %s", sourceURL)
	if err := source.Sync(); err != nil {
		return err
	}
	reg, err := source.Load()
	if err != nil {
		return err
	}
	if pinned {
		a.cfg.RegistryURL = sourceURL
		if err := a.paths.SaveConfig(a.cfg); err != nil {
			return err
		}
		a.source = source
	}
	a.log.Printf("[ok] registry updated: %d tools (schema %s)", len(reg.Tools), reg.Version)
	return nil
}

func (a *App) syncSource(options syncOptions, repo string) (registry.Source, string, bool, error) {
	if options.reference == "" && options.repo == "" {
		return a.source, a.cfg.RegistryURL, false, nil
	}
	selected := options.reference
	if selected == "" {
		selected = "main"
	}
	client, err := registry.NewGitHubClient(repo, a.cfg.RegistryPath, githubToken())
	if err != nil {
		return nil, "", false, err
	}
	a.log.Printf("-> resolving GitHub ref %s from %s", selected, repo)
	commit, err := client.ResolveCommit(context.Background(), selected)
	if err != nil {
		return nil, "", false, err
	}
	sourceURL, err := client.RegistryURL(commit)
	if err != nil {
		return nil, "", false, err
	}
	a.cfg.RegistryRepo = repo
	a.cfg.RegistryRef = selected
	a.cfg.RegistryCommit = commit
	return registry.NewFileSource(a.paths.RegistryFile, sourceURL), sourceURL, true, nil
}

func (a *App) listGitHubRefs(repo string) error {
	client, err := registry.NewGitHubClient(repo, a.cfg.RegistryPath, githubToken())
	if err != nil {
		return err
	}
	refs, err := client.ListTags(context.Background(), 50)
	if err != nil {
		return err
	}
	if len(refs) == 0 {
		a.log.Printf("no GitHub tags found for %s; branches and commit SHAs can still be passed with --ref", repo)
		return nil
	}
	for _, ref := range refs {
		a.log.Printf("%-24s %s", ref.Name, ref.Commit[:12])
	}
	return nil
}

func parseSyncArgs(args []string) (syncOptions, error) {
	var options syncOptions
	for i := 0; i < len(args); i++ {
		argument := args[i]
		switch {
		case argument == "--list-refs":
			options.listRefs = true
		case argument == "--ref" || argument == "--repo":
			if i+1 >= len(args) {
				return syncOptions{}, fmt.Errorf("%s requires a value", argument)
			}
			i++
			if argument == "--ref" {
				options.reference = strings.TrimSpace(args[i])
			} else {
				options.repo = strings.TrimSpace(args[i])
			}
		case strings.HasPrefix(argument, "--ref="):
			options.reference = strings.TrimSpace(strings.TrimPrefix(argument, "--ref="))
		case strings.HasPrefix(argument, "--repo="):
			options.repo = strings.TrimSpace(strings.TrimPrefix(argument, "--repo="))
		default:
			return syncOptions{}, usageError("sync [--list-refs] [--repo owner/repo] [--ref tag|branch|sha]")
		}
	}
	if options.listRefs && options.reference != "" {
		return syncOptions{}, fmt.Errorf("--list-refs and --ref cannot be used together")
	}
	return options, nil
}
