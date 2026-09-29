package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Command describes one package-manager invocation without shell parsing.
type Command struct {
	Name string
	Args []string
	Env  map[string]string
}

// CommandRunner executes installer processes. Tests inject a recording runner.
type CommandRunner interface {
	Run(ctx context.Context, command Command) error
}

type osCommandRunner struct{}

// NewCommandRunner returns the production process runner.
func NewCommandRunner() CommandRunner { return osCommandRunner{} }

func (osCommandRunner) Run(ctx context.Context, command Command) error {
	cmd := exec.CommandContext(ctx, command.Name, command.Args...)
	cmd.Env = mergedEnvironment(command.Env)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", command.Name, err)
	}
	return nil
}

func mergedEnvironment(overrides map[string]string) []string {
	environment := make([]string, 0, len(os.Environ())+len(overrides))
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if _, replaced := overrides[key]; !replaced {
			environment = append(environment, entry)
		}
	}
	for key, value := range overrides {
		environment = append(environment, key+"="+value)
	}
	return environment
}

var _ CommandRunner = osCommandRunner{}
