package cli

import (
	"fmt"

	"github.com/Okymi-X/arsenal/internal/registry"
	"github.com/Okymi-X/arsenal/internal/resolver"
)

func (a *App) cmdUpgrade(args []string) error {
	all := len(args) == 1 && args[0] == "--all"
	if !all && len(args) != 1 {
		return usageError("upgrade <tool>|--all")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	manifest, err := a.store.Load()
	if err != nil {
		return err
	}
	filter := args
	if all {
		filter = nil
	}
	candidates, err := collectUpdateCandidates(reg, manifest, filter)
	if err != nil {
		return err
	}
	upgraded := 0
	for _, candidate := range candidates {
		changed, err := a.upgradeCandidate(reg, candidate, all)
		if err != nil {
			return err
		}
		if changed {
			upgraded++
		}
	}
	if upgraded == 0 {
		a.log.Printf("no upgrades available")
	}
	return nil
}

func (a *App) upgradeCandidate(reg *registry.Registry, candidate updateCandidate, all bool) (bool, error) {
	switch candidate.relation {
	case registry.VersionOutdated:
		return a.upgradeOutdated(reg, candidate)
	case registry.VersionUntracked, registry.VersionNoRecommendation:
		if !all {
			return false, fmt.Errorf("%s@%s has no safe upgrade recommendation", candidate.installed.Name, candidate.installed.Version)
		}
		a.log.Warnf("%s@%s has no safe upgrade recommendation; skipped", candidate.installed.Name, candidate.installed.Version)
	case registry.VersionAhead:
		a.log.Printf("%s@%s is ahead of the tested recommendation; unchanged", candidate.installed.Name, candidate.installed.Version)
	case registry.VersionCurrent:
	}
	return false, nil
}

func (a *App) upgradeOutdated(reg *registry.Registry, candidate updateCandidate) (bool, error) {
	tool, err := reg.MustFindTool(candidate.installed.Name)
	if err != nil {
		return false, err
	}
	currentManifest, err := a.store.Load()
	if err != nil {
		return false, err
	}
	if installed, ok := currentManifest.Find(tool.Name, candidate.recommended.Tag); ok && installationHealthy(installed) {
		if err := a.activateInstalled(currentManifest, installed); err != nil {
			return false, err
		}
		a.log.Printf("[ok] activated existing %s@%s", tool.Name, candidate.recommended.Tag)
		return true, nil
	}
	resolved := resolver.Resolved{Tool: tool, Version: candidate.recommended}
	if err := a.installResolved(resolved); err != nil {
		return false, err
	}
	return true, nil
}
