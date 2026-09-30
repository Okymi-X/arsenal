package cli

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Okymi-X/arsenal/internal/installer"
	"github.com/Okymi-X/arsenal/internal/registry"
	"github.com/Okymi-X/arsenal/internal/resolver"
	"github.com/Okymi-X/arsenal/internal/store"
)

// cmdInstall installs a tool at a resolved version into an isolated backend,
// generates shims, and records the installation in the manifest.
func (a *App) cmdInstall(args []string) error {
	options, err := parseInstallArgs(args)
	if err != nil {
		return err
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	res, err := a.resolveInstallOptions(reg, options)
	if err != nil {
		return err
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

func (a *App) installResolved(res resolver.Resolved) error {
	backend := a.newBackend()
	orch := installer.NewOrchestrator(installer.DefaultMethods(backend, a.paths.Tools))

	a.log.Printf("-> installing %s@%s", res.Tool.Name, res.Version.Tag)
	result, err := orch.Install(a.ctx, res.Tool, res.Version)
	if err != nil {
		return err
	}
	if err := a.linkShims(res.Tool, result.Path); err != nil {
		return err
	}
	if result.Backend == "" {
		result.Backend = a.cfg.DefaultBackend
	}
	if err := a.recordInstall(res, result); err != nil {
		return err
	}
	a.log.Printf("[ok] installed %s@%s", res.Tool.Name, res.Version.Tag)
	return nil
}

// linkShims writes a shim for each binary pointing into the environment and
// makes this the active version.
func (a *App) linkShims(tool registry.Tool, root string) error {
	for _, bin := range tool.AllBinaries() {
		target := filepath.Join(root, "bin", bin)
		if err := a.shims.Write(bin, target); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) recordInstall(res resolver.Resolved, result installer.Result) error {
	m, err := a.store.Load()
	if err != nil {
		return err
	}
	m.Upsert(store.InstalledTool{
		Name:        res.Tool.Name,
		Version:     res.Version.Tag,
		Backend:     result.Backend,
		Path:        result.Path,
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
