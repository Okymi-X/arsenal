package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPathCheckKeepsPrivateShimsIsolated(t *testing.T) {
	binDir := filepath.Join(t.TempDir(), "arsenal", "bin")
	t.Setenv("PATH", strings.Join([]string{"/usr/local/bin", "/usr/bin"}, string(os.PathListSeparator)))

	result := NewPathCheck(binDir).Run()
	if !result.OK || result.Detail != "isolated; use 'arsenal run'" {
		t.Fatalf("Run = %+v", result)
	}
}

func TestPathCheckRejectsGlobalShimExposure(t *testing.T) {
	binDir := filepath.Join(t.TempDir(), "arsenal", "bin")
	t.Setenv("PATH", strings.Join([]string{"/usr/local/bin", binDir, "/usr/bin"}, string(os.PathListSeparator)))

	result := NewPathCheck(binDir).Run()
	if result.OK || !strings.Contains(result.Detail, "may shadow host tools") {
		t.Fatalf("Run = %+v", result)
	}
}
