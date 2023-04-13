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

var DefaultCommander *Commander

func init() {
	DefaultCommander = NewCommander(flag.CommandLine, path.Base(os.Args[0]))
}

func Register(cmd Command, group string) {
	DefaultCommander.Register(cmd, group)
}

func HelpCommand() Command {
	return DefaultCommander.HelpCommand()
}

func FlagsCommand() Command {
	return DefaultCommander.FlagsCommand()
}

func CommandsCommand() Command {
	return DefaultCommander.CommandsCommand()
}

func ImportantFlag(name string) {
	DefaultCommander.ImportantFlag(name)
}

func Execute(ctx context.Context, args ...interface{}) ExitStatus {
	return DefaultCommander.Execute(ctx, args...)
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
