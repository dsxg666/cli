package cli

// An aliaser is a Command wrapping another Command but returning a
// different name as its alias.
type aliaser struct {
	alias string
	Command
}

func (a *aliaser) Name() string { return a.alias }
