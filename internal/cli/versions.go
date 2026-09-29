package cli

import (
	"context"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func (a *App) cmdVersions(args []string) error {
	github := false
	if len(args) == 2 && args[1] == "--github" {
		github = true
	} else if len(args) != 1 {
		return usageError("versions <tool> [--github]")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	tool, err := reg.MustFindTool(args[0])
	if err != nil {
		return err
	}
	if !github {
		for _, version := range tool.Versions {
			status := "untested"
			if version.Tested {
				status = "tested"
			}
			a.log.Printf("%-24s %s", version.Tag, status)
		}
		return nil
	}
	repo, err := registry.ParseGitHubRepo(tool.Repo)
	if err != nil {
		return err
	}
	client, err := registry.NewGitHubClient(repo, "registry.toml", githubToken())
	if err != nil {
		return err
	}
	refs, err := client.ListTags(context.Background(), 50)
	if err != nil {
		return err
	}
	for _, ref := range refs {
		a.log.Printf("%-24s %s", ref.Name, ref.Commit[:12])
	}
	return nil
}
