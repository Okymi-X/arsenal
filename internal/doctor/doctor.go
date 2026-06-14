// Package doctor reports on and repairs broken installs.
//
// It runs a sequence of independent checks, each of which can describe a
// problem and, where possible, fix it. Checks are injected so the command can
// compose them and tests can substitute fakes.
package doctor

// Result is the outcome of a single check.
type Result struct {
	// Name identifies the check.
	Name string
	// OK reports whether the check passed.
	OK bool
	// Detail explains the outcome in one line.
	Detail string
	// Fixable reports whether a failed check can be repaired by Fix.
	Fixable bool
}

// Check is one health check with an optional repair.
type Check interface {
	// Name identifies the check.
	Name() string
	// Run evaluates the check and returns its result.
	Run() Result
	// Fix attempts to repair a failed check.
	Fix() error
}

// Doctor runs a set of checks and optionally repairs failures.
type Doctor struct {
	checks []Check
}

// New returns a Doctor over the given checks.
func New(checks []Check) *Doctor { return &Doctor{checks: checks} }

// Run evaluates every check and returns the results in order.
func (d *Doctor) Run() []Result {
	out := make([]Result, 0, len(d.checks))
	for _, c := range d.checks {
		out = append(out, c.Run())
	}
	return out
}

// Repair runs each failed, fixable check's Fix and returns the names of the
// checks it attempted to repair, along with the first error encountered.
func (d *Doctor) Repair() ([]string, error) {
	var fixed []string
	for _, c := range d.checks {
		res := c.Run()
		if res.OK || !res.Fixable {
			continue
		}
		if err := c.Fix(); err != nil {
			return fixed, err
		}
		fixed = append(fixed, c.Name())
	}
	return fixed, nil
}
