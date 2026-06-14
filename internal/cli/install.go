package cli

import (
	"context"
	"fmt"
	"path/filepath"
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
	if len(args) != 1 {
		return usageError("install <tool>[@version]")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	res, err := resolveSpec(reg, args[0])
	if err != nil {
		return err
	}
	if !res.Version.Tested {
		a.log.Warnf("%s@%s is not marked tested", res.Tool.Name, res.Version.Tag)
	}
	return a.installResolved(reg, res)
}

func (a *App) installResolved(reg *registry.Registry, res resolver.Resolved) error {
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
