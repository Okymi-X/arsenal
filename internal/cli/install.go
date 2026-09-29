package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Okymi-X/arsenal/internal/installer"
	"github.com/Okymi-X/arsenal/internal/isolation"
	"github.com/Okymi-X/arsenal/internal/registry"
	"github.com/Okymi-X/arsenal/internal/resolver"
	"github.com/Okymi-X/arsenal/internal/store"
)

// cmdInstall installs a tool at a resolved version into an isolated backend,
// generates shims, and records the installation in the manifest.
func (a *App) cmdInstall(args []string) error {
	spec, githubRef, err := parseInstallArgs(args)
	if err != nil {
		return err
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	res, err := resolveSpec(reg, spec)
	if err != nil {
		return err
	}
	if githubRef != "" {
		res, err = resolveGitHubInstall(res.Tool, githubRef)
		if err != nil {
			return err
		}
	}
	if !res.Version.Tested {
		a.log.Warnf("%s@%s is not marked tested", res.Tool.Name, res.Version.Tag)
	}
	manifest, err := a.store.Load()
	if err != nil {
		return err
	}
	if installed, ok := manifest.Find(res.Tool.Name, res.Version.Tag); ok && installationHealthy(installed) {
		if err := a.activateInstalled(manifest, installed); err != nil {
			return err
		}
		a.log.Printf("[ok] %s@%s already installed; activated", res.Tool.Name, res.Version.Tag)
		return nil
	}
	return a.installResolved(res)
}

func parseInstallArgs(args []string) (string, string, error) {
	if len(args) == 1 {
		return args[0], "", nil
	}
	if len(args) == 3 && args[1] == "--github-ref" && args[2] != "" {
		if _, _, hasVersion := strings.Cut(args[0], "@"); hasVersion {
			return "", "", fmt.Errorf("catalogue version and --github-ref cannot be combined")
		}
		return args[0], args[2], nil
	}
	return "", "", usageError("install <tool>[@version] [--github-ref tag|branch|sha]")
}

func resolveGitHubInstall(tool registry.Tool, ref string) (resolver.Resolved, error) {
	if tool.InstallMethod != installer.MethodPip && tool.InstallMethod != installer.MethodGitPip {
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
	commit, err := client.ResolveCommit(context.Background(), ref)
	if err != nil {
		return resolver.Resolved{}, err
	}
	tool.InstallMethod = installer.MethodGitPip
	return resolver.Resolved{
		Tool: tool,
		Version: registry.Version{
			Tag:    "github-" + commit[:12],
			Commit: commit,
			Repo:   tool.Repo,
			Tested: false,
			Notes:  "Explicit GitHub ref " + ref,
		},
	}, nil
}

func (a *App) installResolved(res resolver.Resolved) error {
	backend := a.newBackend()
	orch := installer.NewOrchestrator(installer.DefaultMethods(backend))

	a.log.Printf("-> installing %s@%s", res.Tool.Name, res.Version.Tag)
	if err := orch.Install(context.Background(), res.Tool, res.Version); err != nil {
		return err
	}
	if err := a.linkShims(res.Tool, backend); err != nil {
		return err
	}
	if err := a.recordInstall(res, backend); err != nil {
		return err
	}
	a.log.Printf("[ok] installed %s@%s", res.Tool.Name, res.Version.Tag)
	return nil
}

// linkShims writes a shim for each binary pointing into the environment and
// makes this the active version.
func (a *App) linkShims(tool registry.Tool, backend isolation.Backend) error {
	for _, bin := range tool.AllBinaries() {
		target := filepath.Join(backend.Path(), "bin", bin)
		if err := a.shims.Write(bin, target); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) recordInstall(res resolver.Resolved, backend isolation.Backend) error {
	m, err := a.store.Load()
	if err != nil {
		return err
	}
	m.Upsert(store.InstalledTool{
		Name:        res.Tool.Name,
		Version:     res.Version.Tag,
		Backend:     a.cfg.DefaultBackend,
		Path:        backend.Path(),
		Binaries:    res.Tool.AllBinaries(),
		InstalledAt: time.Now().UTC().Format(time.RFC3339),
	})
	if !m.SetActive(res.Tool.Name, res.Version.Tag) {
		return fmt.Errorf("failed to activate %s@%s", res.Tool.Name, res.Version.Tag)
	}
	return a.store.Save(m)
}

// resolveSpec parses and resolves a "tool[@version]" spec against a registry.
func resolveSpec(reg *registry.Registry, spec string) (resolver.Resolved, error) {
	req, err := resolver.ParseRequest(spec)
	if err != nil {
		return resolver.Resolved{}, err
	}
	return resolver.New(reg).Resolve(req)
}
