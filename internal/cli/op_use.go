package cli

import (
	"time"

	"github.com/Okymi-X/arsenal/internal/op"
	"github.com/Okymi-X/arsenal/internal/registry"
)

// opUse resolves an op into a lockfile and installs every pinned tool.
func (a *App) opUse(args []string) error {
	if len(args) != 1 {
		return usageError("op use <name>")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	lf, err := a.lockOp(args[0], reg)
	if err != nil {
		return err
	}
	if err := op.WriteLockfile(a.ops.LockPath(args[0]), lf); err != nil {
		return err
	}
	return a.installLockfile(reg, lf)
}

// opExport resolves an op into a lockfile and writes it for sharing.
func (a *App) opExport(args []string) error {
	if len(args) != 1 {
		return usageError("op export <name>")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	lf, err := a.lockOp(args[0], reg)
	if err != nil {
		return err
	}
	path := a.ops.LockPath(args[0])
	if err := op.WriteLockfile(path, lf); err != nil {
		return err
	}
	a.log.Printf("[ok] exported lockfile -> %s", path)
	return nil
}

// opImport reads a shared lockfile and installs every entry.
func (a *App) opImport(args []string) error {
	if len(args) != 1 {
		return usageError("op import <lockfile>")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	lf, err := op.ReadLockfile(args[0])
	if err != nil {
		return err
	}
	return a.installLockfile(reg, lf)
}

// lockOp loads an op and resolves it into a lockfile against the registry.
func (a *App) lockOp(name string, reg *registry.Registry) (*op.Lockfile, error) {
	o, err := a.ops.Load(name)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	return op.GenerateLockfile(o, reg.Version, now, resolveEntry(reg))
}

// installLockfile installs every entry in a lockfile by tool and version.
func (a *App) installLockfile(reg *registry.Registry, lf *op.Lockfile) error {
	for _, e := range lf.Entries {
		res, err := resolveSpec(reg, e.Tool+"@"+e.Version)
		if err != nil {
			return err
		}
		if err := a.installResolved(res); err != nil {
			return err
		}
	}
	a.log.Printf("[ok] applied lockfile for op %q (%d tools)", lf.Op, len(lf.Entries))
	return nil
}

// resolveEntry builds a ResolveFunc that turns a pin into a lock entry.
func resolveEntry(reg *registry.Registry) op.ResolveFunc {
	return func(tool, version string) (op.LockEntry, error) {
		spec := tool
		if version != "" {
			spec = tool + "@" + version
		}
		res, err := resolveSpec(reg, spec)
		if err != nil {
			return op.LockEntry{}, err
		}
		return op.LockEntry{
			Tool:          res.Tool.Name,
			Version:       res.Version.Tag,
			Commit:        res.Version.Commit,
			PipSpec:       res.Version.PipSpec,
			InstallMethod: res.Tool.InstallMethod,
		}, nil
	}
}
