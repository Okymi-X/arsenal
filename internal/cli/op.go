package cli

import (
	"fmt"
	"time"

	"github.com/Okymi-X/arsenal/internal/op"
)

// cmdOp dispatches the op subcommands.
func (a *App) cmdOp(args []string) error {
	if len(args) == 0 {
		return usageError("op create|use|pin|list|export|import ...")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "create":
		return a.opCreate(rest)
	case "pin":
		return a.opPin(rest)
	case "list":
		return a.opList(rest)
	case "use":
		return a.opUse(rest)
	case "export":
		return a.opExport(rest)
	case "import":
		return a.opImport(rest)
	default:
		return fmt.Errorf("unknown op subcommand %q", sub)
	}
}

func (a *App) opCreate(args []string) error {
	if len(args) < 1 {
		return usageError("op create <name> [description]")
	}
	description := ""
	if len(args) > 1 {
		description = args[1]
	}
	o := &op.Op{
		Name:        args[0],
		Description: description,
		Created:     time.Now().UTC().Format(time.RFC3339),
	}
	if err := a.ops.Create(o); err != nil {
		return err
	}
	a.log.Printf("[ok] created op %q", o.Name)
	return nil
}

func (a *App) opPin(args []string) error {
	if len(args) != 2 {
		return usageError("op pin <name> <tool[@version]>")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	res, err := resolveSpec(reg, args[1])
	if err != nil {
		return err
	}
	o, err := a.ops.Load(args[0])
	if err != nil {
		return err
	}
	o.SetPin(res.Tool.Name, res.Version.Tag)
	if err := a.ops.Save(o); err != nil {
		return err
	}
	a.log.Printf("[ok] pinned %s@%s in op %q", res.Tool.Name, res.Version.Tag, o.Name)
	return nil
}

func (a *App) opList(args []string) error {
	names, err := a.ops.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		a.log.Printf("no ops defined")
		return nil
	}
	for _, n := range names {
		a.log.Printf("%s", n)
	}
	return nil
}
