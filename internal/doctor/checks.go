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

// PathCheck verifies the shim bin directory is present on the user's PATH.
type PathCheck struct {
	binDir string
}

// NewPathCheck returns a PathCheck for the shim bin directory.
func NewPathCheck(binDir string) *PathCheck { return &PathCheck{binDir: binDir} }

// Name identifies the check.
func (c *PathCheck) Name() string { return "shim-path" }

// Run reports whether binDir appears in the PATH environment variable.
func (c *PathCheck) Run() Result {
	for _, p := range filepathList(os.Getenv("PATH")) {
		if p == c.binDir {
			return Result{Name: c.Name(), OK: true, Detail: c.binDir}
		}
	}
	return Result{Name: c.Name(), OK: false, Detail: fmt.Sprintf("add %s to PATH", c.binDir)}
}

// Fix cannot edit the user's shell profile; it reports the manual action.
func (c *PathCheck) Fix() error {
	return fmt.Errorf("add %s to your PATH in your shell profile", c.binDir)
}

func filepathList(path string) []string {
	if path == "" {
		return nil
	}
	return splitList(path)
}
