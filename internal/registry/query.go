package registry

import (
	"fmt"
	"sort"
	"strings"
)

// FindTool returns the tool matching name or one of its aliases.
//
// Matching is case-insensitive. The returned bool is false when no tool
// matches, mirroring the comma-ok idiom for lookups.
func (r *Registry) FindTool(name string) (Tool, bool) {
	want := strings.ToLower(strings.TrimSpace(name))
	for _, t := range r.Tools {
		if strings.ToLower(t.Name) == want {
			return t, true
		}
		for _, a := range t.Aliases {
			if strings.ToLower(a) == want {
				return t, true
			}
		}
	}
	return Tool{}, false
}

// MustFindTool is FindTool that returns an error instead of a bool.
func (r *Registry) MustFindTool(name string) (Tool, error) {
	t, ok := r.FindTool(name)
	if !ok {
		return Tool{}, fmt.Errorf("tool %q not found in registry", name)
	}
	return t, nil
}

// Search returns tools whose name, description, or tags match the query.
//
// Results are ordered by name for stable, predictable output.
func (r *Registry) Search(query string) []Tool {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return r.sortedTools()
	}
	var out []Tool
	for _, t := range r.Tools {
		if toolMatches(t, q) {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// ByCategory returns every tool in the given category, ordered by name.
func (r *Registry) ByCategory(c Category) []Tool {
	var out []Tool
	for _, t := range r.Tools {
		if t.Category == c {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (r *Registry) sortedTools() []Tool {
	out := make([]Tool, len(r.Tools))
	copy(out, r.Tools)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func toolMatches(t Tool, q string) bool {
	if strings.Contains(strings.ToLower(t.Name), q) {
		return true
	}
	if strings.Contains(strings.ToLower(t.Description), q) {
		return true
	}
	for _, tag := range t.Tags {
		if strings.Contains(strings.ToLower(tag), q) {
			return true
		}
	}
	for _, a := range t.Aliases {
		if strings.Contains(strings.ToLower(a), q) {
			return true
		}
	}
	return false
}

// FindVersion returns the version with the given tag.
func (t Tool) FindVersion(tag string) (Version, bool) {
	for _, v := range t.Versions {
		if v.Tag == tag {
			return v, true
		}
	}
	return Version{}, false
}
