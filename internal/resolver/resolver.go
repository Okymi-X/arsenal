// Package resolver turns a user-supplied "tool[@version]" request into a
// concrete tool and version drawn from the registry.
//
// Resolution prefers an explicitly requested version, then the newest tested
// version, and finally the newest catalogued version. The package is pure:
// it depends only on the registry model and has no side effects.
package resolver

import (
	"fmt"
	"strings"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// Request is a parsed "tool[@version]" reference.
type Request struct {
	// Tool is the tool name or alias.
	Tool string
	// Version is the requested tag, or empty to auto-select.
	Version string
}

// Resolved pairs a tool with the single version chosen for it.
type Resolved struct {
	// Tool is the matched catalog entry.
	Tool registry.Tool
	// Version is the selected release.
	Version registry.Version
}

// ParseRequest splits "tool@version" into its components.
func ParseRequest(spec string) (Request, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return Request{}, fmt.Errorf("empty tool specification")
	}
	name, version, found := strings.Cut(spec, "@")
	if found && version == "" {
		return Request{}, fmt.Errorf("version is empty after '@' in %q", spec)
	}
	return Request{Tool: name, Version: version}, nil
}

// Resolver selects a concrete version from the registry for a request.
type Resolver struct {
	reg *registry.Registry
}

// New returns a Resolver backed by the given registry.
func New(reg *registry.Registry) *Resolver { return &Resolver{reg: reg} }

// Resolve matches the request to a tool and a single version.
func (r *Resolver) Resolve(req Request) (Resolved, error) {
	tool, err := r.reg.MustFindTool(req.Tool)
	if err != nil {
		return Resolved{}, err
	}
	v, err := selectVersion(tool, req.Version)
	if err != nil {
		return Resolved{}, err
	}
	return Resolved{Tool: tool, Version: v}, nil
}

func selectVersion(tool registry.Tool, tag string) (registry.Version, error) {
	if tag != "" {
		v, ok := tool.FindVersion(tag)
		if !ok {
			return registry.Version{}, fmt.Errorf("version %q not catalogued for %q", tag, tool.Name)
		}
		return v, nil
	}
	return pickDefault(tool)
}

// pickDefault chooses the newest tested version, else the newest version.
// Versions are stored newest-first, so the first match wins.
func pickDefault(tool registry.Tool) (registry.Version, error) {
	for _, v := range tool.Versions {
		if v.Tested {
			return v, nil
		}
	}
	if len(tool.Versions) > 0 {
		return tool.Versions[0], nil
	}
	return registry.Version{}, fmt.Errorf("tool %q has no versions", tool.Name)
}
