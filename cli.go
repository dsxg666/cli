package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

func init() {
	DefaultCommander = NewCommander(flag.CommandLine, path.Base(os.Args[0]))
}

// Sorting of the commands within a group.
// explainGroup explains all the subcommands for a particular group.
func explainGroup(w io.Writer, group *CommandGroup) {
	if len(group.commands) == 0 {
		return
	}
	if group.name == "" {
		_, _ = fmt.Fprintf(w, "Subcommands:\n")
	} else {
		_, _ = fmt.Fprintf(w, "Subcommands for %s:\n", group.name)
	}
	sort.Sort(group)

	aliases := make(map[string][]string)
	for _, cmd := range group.commands {
		if alias, ok := cmd.(*aliaser); ok {
			root := dealias(alias).Name()

			if _, ok := aliases[root]; !ok {
				aliases[root] = []string{}
			}
			aliases[root] = append(aliases[root], alias.Name())
		}
	}

	for _, cmd := range group.commands {
		if _, ok := cmd.(*aliaser); ok {
			continue
		}

		name := cmd.Name()
		names := []string{name}

		if a, ok := aliases[name]; ok {
			names = append(names, a...)
		}

		_, _ = fmt.Fprintf(w, "\t%-15s  %s\n", strings.Join(names, ", "), cmd.Synopsis())
	}
	_, _ = fmt.Fprintln(w)
}

// explainCmd prints a brief description of a single command.
func explain(w io.Writer, cmd Command) {
	_, _ = fmt.Fprintf(w, "%s", cmd.Usage())
	subflags := flag.NewFlagSet(cmd.Name(), flag.PanicOnError)
	subflags.SetOutput(w)
	cmd.SetFlags(subflags)
	subflags.PrintDefaults()
}

// Alias returns a Command alias which implements a "commands" subcommand.
func Alias(alias string, cmd Command) Command {
	return &aliaser{alias, cmd}
}

// dealias recursivly dealiases a command until a non-aliased command
// is reached.
func dealias(cmd Command) Command {
	if alias, ok := cmd.(*aliaser); ok {
		return dealias(alias.Command)
	}

	return cmd
}

// DefaultCommander is the default commander using flag.CommandLine for flags
// and os.Args[0] for the command name.
var DefaultCommander *Commander

// Register adds a subcommand to the supported subcommands in the
// specified group. (Help output is sorted and arranged by group
// name.)  The empty string is an acceptable group name; such
// subcommands are explained first before named groups. It is a
// wrapper around DefaultCommander.Register.
func Register(cmd Command, group string) {
	DefaultCommander.Register(cmd, group)
}

// ImportantFlag marks a top-level flag as important, which means it
// will be printed out as part of the output of an ordinary "help"
// subcommand.  (All flags, important or not, are printed by the
// "flags" subcommand.) It is a wrapper around
// DefaultCommander.ImportantFlag.
func ImportantFlag(name string) {
	DefaultCommander.ImportantFlag(name)
}

// Execute should be called once the default flags have been
// initialized by flag.Parse. It finds the correct subcommand and
// executes it, and returns an ExitStatus with the result. On a usage
// error, an appropriate message is printed to os.Stderr, and
// ExitUsageError is returned. The additional args are provided as-is
// to the Execute method of the selected Command. It is a wrapper
// around DefaultCommander.Execute.
func Execute(ctx context.Context, args ...interface{}) ExitStatus {
	return DefaultCommander.Execute(ctx, args...)
}

// HelpCommand returns a Command which implements "help" for the
// DefaultCommander. Use Register(HelpCommand(), <group>) for it to be
// recognized.
func HelpCommand() Command {
	return DefaultCommander.HelpCommand()
}

// FlagsCommand returns a Command which implements "flags" for the
// DefaultCommander. Use Register(FlagsCommand(), <group>) for it to be
// recognized.
func FlagsCommand() Command {
	return DefaultCommander.FlagsCommand()
}

// CommandsCommand returns Command which implements a "commands" subcommand.
func CommandsCommand() Command {
	return DefaultCommander.CommandsCommand()
}
