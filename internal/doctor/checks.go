package doctor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Okymi-X/arsenal/internal/config"
)

// DirsCheck verifies the arsenal directory tree exists and can recreate it.
type DirsCheck struct {
	paths config.Paths
}

// NewDirsCheck returns a DirsCheck for the given paths.
func NewDirsCheck(paths config.Paths) *DirsCheck { return &DirsCheck{paths: paths} }

// Name identifies the check.
func (c *DirsCheck) Name() string { return "directories" }

// Run reports whether the root directory exists.
func (c *DirsCheck) Run() Result {
	if _, err := os.Stat(c.paths.Root); err != nil {
		return Result{Name: c.Name(), OK: false, Detail: "root directory missing", Fixable: true}
	}
	return Result{Name: c.Name(), OK: true, Detail: c.paths.Root}
}

// Fix recreates the directory tree.
func (c *DirsCheck) Fix() error { return c.paths.EnsureDirs() }

// PythonCheck verifies the configured Python interpreter is on PATH.
type PythonCheck struct {
	pythonBin string
}

// NewPythonCheck returns a PythonCheck for the given interpreter.
func NewPythonCheck(pythonBin string) *PythonCheck { return &PythonCheck{pythonBin: pythonBin} }

// Name identifies the check.
func (c *PythonCheck) Name() string { return "python" }

// Run reports whether the interpreter resolves on PATH.
func (c *PythonCheck) Run() Result {
	path, err := exec.LookPath(c.pythonBin)
	if err != nil {
		return Result{Name: c.Name(), OK: false, Detail: fmt.Sprintf("%s not found on PATH", c.pythonBin)}
	}
	return Result{Name: c.Name(), OK: true, Detail: path}
}

// Fix cannot install Python; it reports the manual action required.
func (c *PythonCheck) Fix() error {
	return fmt.Errorf("install %s and ensure it is on PATH", c.pythonBin)
}

// PathCheck verifies the private shim directory does not shadow host tools.
type PathCheck struct {
	binDir string
}

// NewPathCheck returns a PathCheck for the shim bin directory.
func NewPathCheck(binDir string) *PathCheck { return &PathCheck{binDir: binDir} }

// Name identifies the check.
func (c *PathCheck) Name() string { return "shim-path" }

// Run reports whether private shims remain isolated from the host PATH.
func (c *PathCheck) Run() Result {
	for _, p := range filepathList(os.Getenv("PATH")) {
		if p == c.binDir {
			return Result{
				Name:   c.Name(),
				OK:     false,
				Detail: fmt.Sprintf("%s is on PATH and may shadow host tools", c.binDir),
			}
		}
	}
	return Result{Name: c.Name(), OK: true, Detail: "isolated; use 'arsenal run'"}
}

// Fix cannot safely infer which shell configuration introduced the path.
func (c *PathCheck) Fix() error {
	return fmt.Errorf("remove %s from PATH to prevent host tool conflicts", c.binDir)
}

func filepathList(path string) []string {
	if path == "" {
		return nil
	}
	return splitList(path)
}
