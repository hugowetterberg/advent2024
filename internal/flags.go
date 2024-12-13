package internal

import (
	"flag"
	"fmt"
	"os"
	"slices"
)

func Args(osArgs []string) ([]string, []string) {
	args := osArgs[1:]

	idx := slices.Index(args, "--")
	if idx == -1 {
		return args, nil
	}

	return args[0:idx], args[idx+1:]
}

func ParseSolutionFlags(set *flag.FlagSet) error {
	_, args := Args(os.Args)

	err := set.Parse(args)
	if err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	return nil
}
