package cli

import (
	"context"
	"flag"
	"fmt"
)

// A flagger is a Command implementing a "flags" command for a given Commander.
type flagger Commander

func (flg *flagger) Name() string           { return "flags" }
func (flg *flagger) Synopsis() string       { return "describe all known top-level flags" }
func (flg *flagger) SetFlags(*flag.FlagSet) {}
func (flg *flagger) Usage() string {
	return `flags [<subcommand>]:
	With an argument, print all flags of <subcommand>. Else,
	print a description of all known top-level flags.  (The basic
	help information only discusses the most generally important
	top-level flags.)
`
}

func (flg *flagger) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) ExitStatus {
	if f.NArg() > 1 {
		f.Usage()
		return ExitUsageError
	}

	if f.NArg() == 0 {
		if flg.topFlags == nil {
			_, _ = fmt.Fprintln(flg.Output, "No top-level flags are defined.")
		} else {
			flg.topFlags.PrintDefaults()
		}
		return ExitSuccess
	}

	for _, group := range flg.commands {
		for _, cmd := range group.commands {
			if f.Arg(0) != cmd.Name() {
				continue
			}
			subflags := flag.NewFlagSet(cmd.Name(), flag.PanicOnError)
			subflags.SetOutput(flg.Output)
			cmd.SetFlags(subflags)
			subflags.PrintDefaults()
			return ExitSuccess
		}
	}
	_, _ = fmt.Fprintf(flg.Error, "Subcommand %s not understood\n", f.Arg(0))
	return ExitFailure
}
