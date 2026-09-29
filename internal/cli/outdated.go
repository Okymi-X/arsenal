package cli

import (
	"fmt"
	"sort"

	"github.com/Okymi-X/arsenal/internal/registry"
	"github.com/Okymi-X/arsenal/internal/store"
)

type updateCandidate struct {
	installed   store.InstalledTool
	recommended registry.Version
	relation    registry.VersionRelation
}

func (a *App) cmdOutdated(args []string) error {
	if len(args) > 1 {
		return usageError("outdated [tool]")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	manifest, err := a.store.Load()
	if err != nil {
		return err
	}
	candidates, err := collectUpdateCandidates(reg, manifest, args)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		a.log.Printf("no active tools installed")
		return nil
	}
	count := 0
	for _, candidate := range candidates {
		switch candidate.relation {
		case registry.VersionOutdated:
			count++
			a.log.Printf("%s %s -> %s", candidate.installed.Name, candidate.installed.Version, candidate.recommended.Tag)
		case registry.VersionAhead:
			a.log.Printf("%s %s is ahead of tested %s", candidate.installed.Name, candidate.installed.Version, candidate.recommended.Tag)
		case registry.VersionUntracked:
			a.log.Warnf("%s@%s is not tracked by the registry", candidate.installed.Name, candidate.installed.Version)
		case registry.VersionNoRecommendation:
			a.log.Warnf("%s has no tested version", candidate.installed.Name)
		case registry.VersionCurrent:
		}
	}
	if count == 0 {
		a.log.Printf("no safe upgrades available")
	}
	return nil
}

func collectUpdateCandidates(reg *registry.Registry, manifest *store.Manifest, args []string) ([]updateCandidate, error) {
	requested := ""
	if len(args) == 1 {
		tool, err := reg.MustFindTool(args[0])
		if err != nil {
			return nil, err
		}
		requested = tool.Name
	}
	var candidates []updateCandidate
	for _, installed := range manifest.Tools {
		if !installed.Active || requested != "" && installed.Name != requested {
			continue
		}
		tool, ok := reg.FindTool(installed.Name)
		if !ok {
			candidates = append(candidates, updateCandidate{installed: installed, relation: registry.VersionUntracked})
			continue
		}
		recommended, relation := tool.CompareToRecommended(installed.Version)
		candidates = append(candidates, updateCandidate{installed: installed, recommended: recommended, relation: relation})
	}
	if requested != "" && len(candidates) == 0 {
		return nil, fmt.Errorf("%s has no active installation", requested)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].installed.Name < candidates[j].installed.Name })
	return candidates, nil
}
