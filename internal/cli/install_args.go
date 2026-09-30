package cli

import (
	"fmt"
	"strings"
)

type installOptions struct {
	spec          string
	query         string
	githubRef     string
	selectTool    bool
	selectVersion bool
	selectGitHub  bool
}

func parseInstallArgs(args []string) (installOptions, error) {
	if len(args) == 1 {
		if args[0] == "--select" {
			return installOptions{selectTool: true, selectVersion: true}, nil
		}
		return installOptions{spec: args[0]}, nil
	}
	if len(args) == 2 && args[0] == "--select" {
		return installOptions{query: args[1], selectTool: true, selectVersion: true}, nil
	}
	if len(args) == 2 && args[1] == "--select" {
		if _, _, hasVersion := strings.Cut(args[0], "@"); hasVersion {
			return installOptions{}, fmt.Errorf("catalogue version and --select cannot be combined")
		}
		return installOptions{spec: args[0], selectVersion: true}, nil
	}
	if len(args) == 2 && args[1] == "--github-select" {
		if _, _, hasVersion := strings.Cut(args[0], "@"); hasVersion {
			return installOptions{}, fmt.Errorf("catalogue version and --github-select cannot be combined")
		}
		return installOptions{spec: args[0], selectGitHub: true}, nil
	}
	if len(args) == 3 && args[1] == "--github-ref" && args[2] != "" {
		if _, _, hasVersion := strings.Cut(args[0], "@"); hasVersion {
			return installOptions{}, fmt.Errorf("catalogue version and --github-ref cannot be combined")
		}
		return installOptions{spec: args[0], githubRef: args[2]}, nil
	}
	return installOptions{}, usageError("install <tool>[@version] [--select|--github-select|--github-ref tag|branch|sha] | install --select [query]")
}
