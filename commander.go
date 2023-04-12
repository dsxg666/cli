package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

// A Commander represents a set of commands.
type Commander struct {
	commands  []*CommandGroup
	topFlags  *flag.FlagSet // top-level flags
	important []string      // important top-level flags
	name      string        // normally path.Base(os.Args[0])

	Explain        func(io.Writer)                // A function to print a top level usage explanation. Can be overridden.
	ExplainGroup   func(io.Writer, *CommandGroup) // A function to print a command group's usage explanation. Can be overridden.
	ExplainCommand func(io.Writer, Command)       // A function to print a command usage explanation. Can be overridden.

	Output io.Writer // Output specifies where the commander should write its output (default: os.Stdout).
	Error  io.Writer // Error specifies where the commander should write its error (default: os.Stderr).
}

// NewCommander returns a new commander with the specified top-level
// flags and command name. The Usage function for the topLevelFlags
// will be set as well.
func NewCommander(topLevelFlags *flag.FlagSet, name string) *Commander {
	cdr := &Commander{
		topFlags: topLevelFlags,
		name:     name,
		Output:   os.Stdout,
		Error:    os.Stderr,
	}

	cdr.Explain = cdr.explain
	cdr.ExplainGroup = explainGroup
	cdr.ExplainCommand = explain
	topLevelFlags.Usage = func() { cdr.Explain(cdr.Error) }
	return cdr
}

// Name returns the commander's name
func (cdr *Commander) Name() string {
	return cdr.name
}

// Register adds a subcommand to the supported subcommands in the
// specified group. (Help output is sorted and arranged by group name.)
// The empty string is an acceptable group name; such subcommands are
// explained first before named groups.
func (cdr *Commander) Register(cmd Command, group string) {
	for _, g := range cdr.commands {
		if g.name == group {
			g.commands = append(g.commands, cmd)
			return
		}
	}
	cdr.commands = append(cdr.commands, &CommandGroup{
		name:     group,
		commands: []Command{cmd},
	})
}

// ImportantFlag marks a top-level flag as important, which means it
// will be printed out as part of the output of an ordinary "help"
// subcommand.  (All flags, important or not, are printed by the
// "flags" subcommand.)
func (cdr *Commander) ImportantFlag(name string) {
	cdr.important = append(cdr.important, name)
}

// VisitGroups visits each command group in lexicographical order, calling
// fn for each.
func (cdr *Commander) VisitGroups(fn func(*CommandGroup)) {
	sort.Sort(byGroupName(cdr.commands))
	for _, g := range cdr.commands {
		fn(g)
	}
}

// VisitCommands visits each command in registered order grouped by
// command group in lexicographical order, calling fn for each.
func (cdr *Commander) VisitCommands(fn func(*CommandGroup, Command)) {
	cdr.VisitGroups(func(g *CommandGroup) {
		for _, cmd := range g.commands {
			fn(g, cmd)
		}
	})
}

// VisitAllImportant visits the important top level flags in lexicographical
// order, calling fn for each. It visits all flags, even those not set.
func (cdr *Commander) VisitAllImportant(fn func(*flag.Flag)) {
	sort.Strings(cdr.important)
	for _, name := range cdr.important {
		f := cdr.topFlags.Lookup(name)
		if f == nil {
			panic(fmt.Sprintf("Important flag (%s) is not defined", name))
		}
		fn(f)
	}
}

// VisitAll visits the top level flags in lexicographical order, calling fn
// for each. It visits all flags, even those not set.
func (cdr *Commander) VisitAll(fn func(*flag.Flag)) {
	if cdr.topFlags != nil {
		cdr.topFlags.VisitAll(fn)
	}
}

// countFlags returns the number of top-level flags defined, even those not set.
func (cdr *Commander) countTopFlags() int {
	count := 0
	cdr.VisitAll(func(*flag.Flag) {
		count++
	})
	return count
}

// Execute should be called once the top-level-flags on a Commander
// have been initialized. It finds the correct subcommand and executes
// it, and returns an ExitStatus with the result. On a usage error, an
// appropriate message is printed to os.Stderr, and ExitUsageError is
// returned. The additional args are provided as-is to the Execute method
// of the selected Command.
func (cdr *Commander) Execute(ctx context.Context, args ...interface{}) ExitStatus {
	if cdr.topFlags.NArg() < 1 {
		cdr.topFlags.Usage()
		return ExitUsageError
	}

	name := cdr.topFlags.Arg(0)

	for _, group := range cdr.commands {
		for _, cmd := range group.commands {
			if name != cmd.Name() {
				continue
			}
			f := flag.NewFlagSet(name, flag.ContinueOnError)
			f.Usage = func() { cdr.ExplainCommand(cdr.Error, cmd) }
			cmd.SetFlags(f)
			if f.Parse(cdr.topFlags.Args()[1:]) != nil {
				return ExitUsageError
			}
			return cmd.Execute(ctx, f, args...)
		}
	}

	// Cannot find this command.
	cdr.topFlags.Usage()
	return ExitUsageError
}

// explain prints a brief description of all the subcommands and the
// important top-level flags.
func (cdr *Commander) explain(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Usage: %s <flags> <subcommand> <subcommand args>\n\n", cdr.name)
	sort.Sort(byGroupName(cdr.commands))
	for _, group := range cdr.commands {
		cdr.ExplainGroup(w, group)
	}
	if cdr.topFlags == nil {
		_, _ = fmt.Fprintln(w, "\nNo top level flags.")
		return
	}

	sort.Strings(cdr.important)
	if len(cdr.important) == 0 {
		if cdr.countTopFlags() > 0 {
			_, _ = fmt.Fprintf(w, "\nUse \"%s flags\" for a list of top-level flags\n", cdr.name)
		}
		return
	}

	_, _ = fmt.Fprintf(w, "\nTop-level flags (use \"%s flags\" for a full list):\n", cdr.name)
	for _, name := range cdr.important {
		f := cdr.topFlags.Lookup(name)
		if f == nil {
			panic(fmt.Sprintf("Important flag (%s) is not defined", name))
		}
		_, _ = fmt.Fprintf(w, "  -%s=%s: %s\n", f.Name, f.DefValue, f.Usage)
	}
}

// HelpCommand returns a Command which implements a "help" subcommand.
func (cdr *Commander) HelpCommand() Command {
	return (*helper)(cdr)
}

// FlagsCommand returns a Command which implements a "flags" subcommand.
func (cdr *Commander) FlagsCommand() Command {
	return (*flagger)(cdr)
}

// CommandsCommand returns Command which implements a "commands" subcommand.
func (cdr *Commander) CommandsCommand() Command {
	return (*lister)(cdr)
}
