package args

import (
	"strings"
	"testing"

	"github.com/shabbyrobe/cmdy/internal/assert"
)

func TestArgSetOneString(t *testing.T) {
	type vals struct {
		foo string
	}
	setup := func() (*ArgSet, *vals) {
		var v vals
		as := NewArgSet()
		as.String(&v.foo, "foo", "Usage...")
		return as, &v
	}

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		err := as.Parse([]string{})
		tt.MustEqual("", v.foo)
		tt.MustAssert(strings.Contains(err.Error(), "missing arg <foo> at position 1"), err.Error())
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		tt.MustOK(as.Parse([]string{"a"}))
		tt.MustEqual("a", v.foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		err := as.Parse([]string{"a", "b"})
		tt.MustEqual("a", v.foo)
		tt.MustAssert(strings.Contains(err.Error(), "found 1 additional arg"), err.Error())
	})
}

func TestNonOptionalArgAfterOptionalArg(t *testing.T) {
	type vals struct {
		foo, bar, baz string
	}
	setup := func() (*ArgSet, *vals) {
		var v vals
		as := NewArgSet()
		as.String(&v.foo, "foo", "Usage...")
		as.StringOptional(&v.bar, "bar", "default", "Usage...")
		as.String(&v.baz, "baz", "Usage...")
		return as, &v
	}

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		tt.MustOK(as.Parse([]string{"a", "b", "c"}))
		tt.MustEqual("a", v.foo)
		tt.MustEqual("b", v.bar)
		tt.MustEqual("c", v.baz)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		tt.MustOK(as.Parse([]string{"a"}))
		tt.MustEqual("a", v.foo)
		tt.MustEqual("default", v.bar)
		tt.MustEqual("", v.baz)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		err := as.Parse([]string{})
		tt.MustAssert(err != nil) // FIXME: check error
		tt.MustEqual("", v.foo)
		tt.MustEqual("default", v.bar)
		tt.MustEqual("", v.baz)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		err := as.Parse([]string{"a", "b", "c", "d"})
		tt.MustAssert(err != nil) // FIXME: check error
		tt.MustEqual("a", v.foo)
		tt.MustEqual("b", v.bar)
		tt.MustEqual("c", v.baz)
	})
}

func TestRemainingOnly(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", AnyLen, "Usage...")

		tt.MustOK(as.Parse([]string{"a", "b", "c"}))
		tt.MustEqual([]string{"a", "b", "c"}, foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", AnyLen, "Usage...")

		tt.MustOK(as.Parse([]string{}))
		tt.MustEqual(0, len(foo))
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", Min(1), "Usage...")

		err := as.Parse([]string{})
		tt.MustAssert(err != nil) // FIXME: check error
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", Min(1), "Usage...")

		tt.MustOK(as.Parse([]string{"a"}))
		tt.MustEqual([]string{"a"}, foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", Max(1), "Usage...")

		tt.MustOK(as.Parse([]string{"a"}))
		tt.MustEqual([]string{"a"}, foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", Max(1), "Usage...")

		err := as.Parse([]string{"a", "b"})
		tt.MustAssert(err != nil) // FIXME: check error
	})
}

func TestRemainingInts(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []int
		as := NewArgSet()
		as.RemainingInts(&foo, "foo", AnyLen, "Usage...")

		tt.MustOK(as.Parse([]string{"1", "2", "3"}))
		tt.MustEqual([]int{1, 2, 3}, foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []int
		as := NewArgSet()
		as.RemainingInts(&foo, "foo", AnyLen, "Usage...")

		err := as.Parse([]string{"1", "2", "quack"})
		tt.MustAssert(err != nil) // FIXME: check error
	})
}

func TestRemainingAfterOptional(t *testing.T) {
	type vals struct {
		foo string
		rem []string
	}
	setup := func() (*ArgSet, *vals) {
		var v vals
		as := NewArgSet()
		as.StringOptional(&v.foo, "foo", "default", "Usage...")
		as.Remaining(&v.rem, "rem", AnyLen, "Usage...")
		return as, &v
	}

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		tt.MustOK(as.Parse([]string{"a", "b", "c"}))
		tt.MustEqual("a", v.foo)
		tt.MustEqual([]string{"b", "c"}, v.rem)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		tt.MustOK(as.Parse([]string{"a"}))
		tt.MustEqual("a", v.foo)
		tt.MustEqual(0, len(v.rem))
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		tt.MustOK(as.Parse([]string{}))
		tt.MustEqual("default", v.foo)
		tt.MustEqual(0, len(v.rem))
	})
}

type (
	hintOnlyVar string
	kindOnlyVar string
	hintBothVar string
	hintNoneVar string
)

func (h hintOnlyVar) String() string            { return string(h) }
func (h hintOnlyVar) Hint() (kind, hint string) { return "", "hint" }
func (h *hintOnlyVar) Set(s string) error       { *h = hintOnlyVar(s); return nil }

func (h kindOnlyVar) String() string            { return string(h) }
func (h kindOnlyVar) Hint() (kind, hint string) { return "kind", "" }
func (h *kindOnlyVar) Set(s string) error       { *h = kindOnlyVar(s); return nil }

func (h hintBothVar) String() string            { return string(h) }
func (h hintBothVar) Hint() (kind, hint string) { return "kind", "hint" }
func (h *hintBothVar) Set(s string) error       { *h = hintBothVar(s); return nil }

func (h hintNoneVar) String() string            { return string(h) }
func (h hintNoneVar) Hint() (kind, hint string) { return "", "" }
func (h *hintNoneVar) Set(s string) error       { *h = hintNoneVar(s); return nil }

const expectedHintableUsage = `
  <hintonly> hint
        hint only
  <kindonly> (kind)
        kind only
  <hintboth> (kind) hint
        hint both
  <hintnone>
        hint none
`

func TestArgSetHintable(t *testing.T) {
	tt := assert.WrapTB(t)

	var hintOnly hintOnlyVar
	var kindOnly kindOnlyVar
	var hintBoth hintBothVar
	var hintNone hintNoneVar

	set := NewArgSet()
	set.Var(&hintOnly, "hintonly", "hint only")
	set.Var(&kindOnly, "kindonly", "kind only")
	set.Var(&hintBoth, "hintboth", "hint both")
	set.Var(&hintNone, "hintnone", "hint none")

	// FIXME: brittle test, but adequate for now.
	tt.MustEqual(expectedHintableUsage, "\n"+set.Usage())
}
