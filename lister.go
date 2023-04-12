package main

import (
	"context"
	"flag"
	"fmt"
)

// A lister is a Command implementing a "commands" command for a given Commander.
type lister Commander

func (l *lister) Name() string           { return "commands" }
func (l *lister) Synopsis() string       { return "list all command names" }
func (l *lister) SetFlags(*flag.FlagSet) {}
func (l *lister) Usage() string {
	return `commands:
	Print a list of all commands.
`
}
func (l *lister) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) ExitStatus {
	if f.NArg() != 0 {
		f.Usage()
		return ExitUsageError
	}

	for _, group := range l.commands {
		for _, cmd := range group.commands {
			fmt.Fprintf(l.Output, "%s\n", cmd.Name())
		}
	}
	return ExitSuccess
}
