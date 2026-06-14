package op

import "fmt"

// ResolveFunc resolves a single pin to a fully pinned lock entry.
//
// It is injected by the caller (wired to the registry and resolver) so the op
// package stays free of registry dependencies.
type ResolveFunc func(tool, version string) (LockEntry, error)

// GenerateLockfile resolves every pin in an op into a lockfile.
func GenerateLockfile(o *Op, registryVersion, now string, resolve ResolveFunc) (*Lockfile, error) {
	lf := &Lockfile{
		Op:              o.Name,
		Generated:       now,
		RegistryVersion: registryVersion,
	}
	for _, p := range o.Pins {
		entry, err := resolve(p.Tool, p.Version)
		if err != nil {
			return nil, fmt.Errorf("resolve pin %s@%s: %w", p.Tool, p.Version, err)
		}
		lf.Entries = append(lf.Entries, entry)
	}
	if err := validateLockfile(lf); err != nil {
		return nil, err
	}
	return lf, nil
}
