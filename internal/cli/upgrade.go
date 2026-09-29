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
		switch candidate.relation {
		case registry.VersionOutdated:
			tool, _ := reg.FindTool(candidate.installed.Name)
			resolved := resolver.Resolved{Tool: tool, Version: candidate.recommended}
			currentManifest, err := a.store.Load()
			if err != nil {
				return err
			}
			if installed, ok := currentManifest.Find(tool.Name, candidate.recommended.Tag); ok && installationHealthy(installed) {
				if err := a.activateInstalled(currentManifest, installed); err != nil {
					return err
				}
				a.log.Printf("[ok] activated existing %s@%s", tool.Name, candidate.recommended.Tag)
			} else if err := a.installResolved(resolved); err != nil {
				return err
			}
			upgraded++
		case registry.VersionUntracked, registry.VersionNoRecommendation:
			message := fmt.Sprintf("%s@%s has no safe upgrade recommendation", candidate.installed.Name, candidate.installed.Version)
			if !all {
				return fmt.Errorf("%s", message)
			}
			a.log.Warnf("%s; skipped", message)
		case registry.VersionAhead:
			a.log.Printf("%s@%s is ahead of the tested recommendation; unchanged", candidate.installed.Name, candidate.installed.Version)
		case registry.VersionCurrent:
		}
	}
	if upgraded == 0 {
		a.log.Printf("no upgrades available")
	}
	return nil
}
