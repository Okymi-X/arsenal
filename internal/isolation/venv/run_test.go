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

	err := runCommand(ctx, os.Args[0], "-test.run=^$")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("runCommand error = %v, want context.Canceled", err)
	}
}
