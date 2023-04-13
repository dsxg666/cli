# cli
A Go Command-line Interface.

## Usage
```go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dsxg666/cli"
)

type printCmd struct {
	capitalize bool
}

func (*printCmd) Name() string     { return "print" }
func (*printCmd) Synopsis() string { return "Print args to stdout." }
func (*printCmd) Usage() string {
	return `print [-capitalize] <some text>:
  Print args to stdout.
`
}

func (p *printCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.capitalize, "capitalize", false, "capitalize output")
}

func (p *printCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) cli.ExitStatus {
	for _, arg := range f.Args() {
		if p.capitalize {
			arg = strings.ToUpper(arg)
		}
		fmt.Printf("%s ", arg)
	}
	fmt.Println()
	return cli.ExitSuccess
}

func main() {
	cli.Register(cli.HelpCommand(), "")
	cli.Register(cli.FlagsCommand(), "")
	cli.Register(cli.CommandsCommand(), "")
	cli.Register(&printCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(cli.Execute(ctx)))
}

```