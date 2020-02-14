package cmdy

import (
	"bytes"

	"github.com/ArtProcessors/cmdy/arg"
)

type testCmd struct {
	synopsis  string
	usage     string
	flags     *FlagSet
	args      *arg.ArgSet
	configure func(flags *FlagSet, args *arg.ArgSet)

	err error // Takes precedence over 'run'
	run func(c Context) error
}

func testCmdRunBuilder(r func(c Context) error) func() Command {
	return func() Command { return &testCmd{run: r} }
}

func testBuilder(c Command) func() Command {
	return func() Command { return c }
}

func (t *testCmd) Synopsis() string  { return t.synopsis }
func (t *testCmd) Usage() string     { return t.usage }
func (t *testCmd) Flags() *FlagSet   { return t.flags }
func (t *testCmd) Args() *arg.ArgSet { return t.args }

func (t *testCmd) Configure(flags *FlagSet, args *arg.ArgSet) {
	if t.configure != nil {
		t.configure(flags, args)
	}
}

func (t *testCmd) Run(c Context) error {
	if t.err != nil {
		return t.err
	}
	if t.run != nil {
		return t.run(c)
	}
	return nil
}

type testRunner struct {
	stdin  bytes.Buffer
	stdout bytes.Buffer
	stderr bytes.Buffer
	*Runner
}

func newTestRunner() *testRunner {
	tr := &testRunner{}
	tr.Runner = &Runner{
		Stdin:  &tr.stdin,
		Stdout: &tr.stdout,
		Stderr: &tr.stderr,
	}
	return tr
}
