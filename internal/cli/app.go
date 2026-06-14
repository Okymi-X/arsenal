package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/Okymi-X/arsenal/internal/config"
	"github.com/Okymi-X/arsenal/internal/isolation"
	"github.com/Okymi-X/arsenal/internal/isolation/container"
	"github.com/Okymi-X/arsenal/internal/isolation/venv"
	"github.com/Okymi-X/arsenal/internal/logx"
	"github.com/Okymi-X/arsenal/internal/op"
	"github.com/Okymi-X/arsenal/internal/registry"
	"github.com/Okymi-X/arsenal/internal/shim"
	"github.com/Okymi-X/arsenal/internal/store"
	builtin "github.com/Okymi-X/arsenal/registry"
)

// Options carries the wiring inputs constructed by main.
type Options struct {
	// Paths are the resolved filesystem locations.
	Paths config.Paths
	// Cfg is the loaded user configuration.
	Cfg config.Config
	// Version is the build version string.
	Version string
	// Stdout receives primary command output.
	Stdout io.Writer
	// Stderr receives diagnostics.
	Stderr io.Writer
}

// App holds the wired dependencies shared by all commands.
//
// It is constructed once per process and passed explicitly to each command;
// there is no package-level mutable state.
type App struct {
	paths   config.Paths
	cfg     config.Config
	version string
	out     io.Writer
	errw    io.Writer
	log     *logx.Logger
	source  registry.Source
	store   store.Store
	shims   *shim.Manager
	ops     *op.Manager
}

// New constructs an App from wiring options.
func New(opts Options) *App {
	return &App{
		paths:   opts.Paths,
		cfg:     opts.Cfg,
		version: opts.Version,
		out:     opts.Stdout,
		errw:    opts.Stderr,
		log:     logx.New(logx.LevelQuiet, opts.Stdout, opts.Stderr),
		source:  registry.NewFileSource(opts.Paths.RegistryFile, opts.Cfg.RegistryURL),
		store:   store.NewFileStore(opts.Paths.ManifestFile),
		shims:   shim.NewManager(opts.Paths.Bin),
		ops:     op.NewManager(opts.Paths.Ops),
	}
}

// setVerbose rebuilds the logger at the requested verbosity.
func (a *App) setVerbose(v bool) {
	level := logx.LevelQuiet
	if v {
		level = logx.LevelVerbose
	}
	a.log = logx.New(level, a.out, a.errw)
}

// loadRegistry returns the active registry, seeding the embedded copy on first
// use so the tool works offline out of the box.
func (a *App) loadRegistry() (*registry.Registry, error) {
	if _, err := os.Stat(a.paths.RegistryFile); err != nil {
		if err := a.seedRegistry(); err != nil {
			return nil, err
		}
	}
	return a.source.Load()
}

func (a *App) seedRegistry() error {
	if err := a.paths.EnsureDirs(); err != nil {
		return err
	}
	if err := os.WriteFile(a.paths.RegistryFile, builtin.Bytes(), 0o644); err != nil {
		return fmt.Errorf("seed registry: %w", err)
	}
	return nil
}

// canonicalName resolves a tool name or alias to its canonical registry name.
// If the registry cannot be loaded or the tool is unknown, the input is
// returned unchanged so commands still work against the manifest directly.
func (a *App) canonicalName(name string) string {
	reg, err := a.loadRegistry()
	if err != nil {
		return name
	}
	if tool, ok := reg.FindTool(name); ok {
		return tool.Name
	}
	return name
}

// newBackend constructs the isolation backend selected by configuration.
func (a *App) newBackend() isolation.Backend {
	if a.cfg.DefaultBackend == "container" {
		return container.New("podman")
	}
	return venv.New(a.cfg.PythonBin, a.paths.Tools)
}
