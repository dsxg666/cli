package cli

import (
	"context"
	"flag"
)

// A Command represents a single command.
type Command interface {
	// Name returns the name of the command.
	Name() string

	// Synopsis returns a short string (less than one line) describing the command.
	Synopsis() string

	// Usage returns a long string explaining the command and giving usage information.
	Usage() string

	// SetFlags adds the flags for this command to the specified set.
	SetFlags(*flag.FlagSet)

	// Execute executes the command and returns an ExitStatus.
	Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) ExitStatus
}
