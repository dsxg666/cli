package cli

import (
	"context"
	"flag"
	"fmt"
)

type helper Commander

func (h *helper) Name() string           { return "help" }
func (h *helper) Synopsis() string       { return "describe subcommands and their syntax" }
func (h *helper) SetFlags(*flag.FlagSet) {}
func (h *helper) Usage() string {
	return `help [<subcommand>]:
	With an argument, prints detailed information on the use of
	the specified subcommand. With no argument, print a list of
	all commands and a brief description of each.
`
}

func (h *helper) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) ExitStatus {
	switch f.NArg() {
	case 0:
		(*Commander)(h).Explain(h.Output)
		return ExitSuccess

	case 1:
		for _, group := range h.commands {
			for _, cmd := range group.commands {
				if f.Arg(0) != cmd.Name() {
					continue
				}
				(*Commander)(h).ExplainCommand(h.Output, cmd)
				return ExitSuccess
			}
		}
		_, _ = fmt.Fprintf(h.Error, "Subcommand %s not understood\n", f.Arg(0))
	}

	f.Usage()
	return ExitUsageError
}
