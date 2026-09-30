package venv

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestRunCommandHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := runCommandEnv(ctx, os.Environ(), os.Args[0], "-test.run=^$")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("runCommandEnv error = %v, want context.Canceled", err)
	}
}

func TestRunRejectsUnsafeBinaryName(t *testing.T) {
	backend := &Backend{dir: t.TempDir()}
	if err := backend.Run(context.Background(), []string{"../host-tool"}); err == nil {
		t.Fatal("Run accepted a binary outside the environment")
	}
}
