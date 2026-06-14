package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Okymi-X/arsenal/internal/fetcher"
	"github.com/Okymi-X/arsenal/internal/registry"
)

// cmdFetch pulls the latest version of a precompiled binary (a release asset or
// a raw repository file) into a destination directory for staging onto a
// target. It creates no isolated environment and no shim.
func (a *App) cmdFetch(args []string) error {
	opts, err := parseFetchArgs(args)
	if err != nil {
		return err
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	asset, err := reg.MustFindAsset(opts.name)
	if err != nil {
		return err
	}

	f := fetcher.New(githubToken())
	ctx := context.Background()
	if opts.list {
		return a.listAsset(ctx, f, asset, opts.build)
	}

	a.log.Infof("-> fetching %s", asset.Name)
	res, err := f.Fetch(ctx, asset, fetcher.Selection{
		Binary:  opts.binary,
		Build:   opts.build,
		DestDir: opts.dest,
	})
	if err != nil {
		return err
	}
	a.log.Printf("[ok] fetched %s %s@%s -> %s (%s)",
		asset.Name, res.File, res.Version, res.Path, humanSize(res.Size))
	return nil
}

func (a *App) listAsset(ctx context.Context, f *fetcher.Fetcher, asset registry.Asset, build string) error {
	names, err := f.List(ctx, asset, build)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		a.log.Printf("no files available for %s", asset.Name)
		return nil
	}
	for _, n := range names {
		a.log.Printf("%s", n)
	}
	return nil
}

type fetchOpts struct {
	name   string
	binary string
	build  string
	dest   string
	list   bool
}

// parseFetchArgs parses "fetch <asset> [binary] [--dest D] [--build B] [--list]".
// Flags accept either "--flag value" or "--flag=value".
func parseFetchArgs(args []string) (fetchOpts, error) {
	o := fetchOpts{dest: "."}
	var pos []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			pos = append(pos, arg)
			continue
		}
		key, val, inline := splitFlag(arg)
		switch key {
		case "--list":
			o.list = true
		case "--dest", "--build":
			if !inline {
				i++
				if i >= len(args) {
					return o, fmt.Errorf("flag %s needs a value", key)
				}
				val = args[i]
			}
			if key == "--dest" {
				o.dest = val
			} else {
				o.build = val
			}
		default:
			return o, fmt.Errorf("unknown flag %q", arg)
		}
	}
	if len(pos) == 0 || len(pos) > 2 {
		return o, usageError("fetch <asset> [binary] [--dest <dir>] [--build <build>] [--list]")
	}
	o.name = pos[0]
	if len(pos) == 2 {
		o.binary = pos[1]
	}
	return o, nil
}

// splitFlag splits "--flag=value" into key and value; inline reports whether a
// value was attached with "=".
func splitFlag(arg string) (key, val string, inline bool) {
	if eq := strings.IndexByte(arg, '='); eq >= 0 {
		return arg[:eq], arg[eq+1:], true
	}
	return arg, "", false
}

func githubToken() string {
	for _, k := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}
