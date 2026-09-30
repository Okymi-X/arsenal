package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Okymi-X/arsenal/internal/pipspec"
	"github.com/Okymi-X/arsenal/internal/registry"
)

// checker verifies registry versions against their official sources.
type checker struct {
	http        *http.Client
	githubToken string
	pypiBase    string
}

// newChecker builds a checker with the given HTTP timeout and optional GitHub
// token (used to raise the API rate limit in CI).
func newChecker(timeout time.Duration, githubToken string) *checker {
	return &checker{
		http:        &http.Client{Timeout: timeout},
		githubToken: githubToken,
		pypiBase:    "https://pypi.org",
	}
}

// verify checks that a single version exists upstream for its install method.
func (c *checker) verify(tool registry.Tool, v registry.Version) error {
	if err := c.verifyPipDependencies(v.PipDependencies); err != nil {
		return err
	}
	switch tool.InstallMethod {
	case "pip":
		pkg, ver, err := pipTarget(tool, v)
		if err != nil {
			return err
		}
		return c.checkPyPI(pkg, ver)
	case "gitpip", "gobin", "cargo", "binary":
		return c.checkGitRef(tool, v)
	default:
		return fmt.Errorf("unknown install method %q", tool.InstallMethod)
	}
}

func (c *checker) verifyPipDependencies(requirements []string) error {
	for _, requirement := range requirements {
		pin, err := pipspec.ParseExact(requirement)
		if err != nil {
			return err
		}
		if err := c.checkPyPI(pin.Name, pin.Version); err != nil {
			return fmt.Errorf("verify dependency %s: %w", requirement, err)
		}
	}
	return nil
}

// pipTarget derives the PyPI package name and version from a version entry,
// preferring the explicit pip_spec.
func pipTarget(tool registry.Tool, v registry.Version) (pkg, ver string, err error) {
	if spec := strings.TrimSpace(v.PipSpec); spec != "" {
		pin, parseErr := pipspec.ParseExact(spec)
		if parseErr != nil {
			return "", "", parseErr
		}
		return pin.Name, pin.Version, nil
	}
	return tool.Name, v.Tag, nil
}

// refCandidates returns the Git refs to try for a non-pip tool: the pinned
// commit, the tag, and the tag with a leading "v".
func refCandidates(v registry.Version) []string {
	var refs []string
	if c := strings.TrimSpace(v.Commit); c != "" {
		refs = append(refs, c)
	}
	if v.Tag != "" {
		refs = append(refs, v.Tag, "v"+v.Tag)
	}
	return refs
}
