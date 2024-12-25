package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/hugowetterberg/advent2024/d01"
	"github.com/hugowetterberg/advent2024/d02"
	"github.com/hugowetterberg/advent2024/d03"
	"github.com/hugowetterberg/advent2024/d04"
	"github.com/hugowetterberg/advent2024/d05"
	"github.com/hugowetterberg/advent2024/d06"
	"github.com/hugowetterberg/advent2024/d07"
	"github.com/hugowetterberg/advent2024/d08"
	"github.com/hugowetterberg/advent2024/d09"
	"github.com/hugowetterberg/advent2024/d10"
	"github.com/hugowetterberg/advent2024/d11"
	"github.com/hugowetterberg/advent2024/d12"
	"github.com/hugowetterberg/advent2024/internal"
)

type SolutionFunc func(input io.Reader) error

var days = [][]SolutionFunc{
	{d01.SolutionOne, d01.SolutionTwo},
	{d02.SolutionOne, d02.SolutionTwo},
	{d03.SolutionOne, d03.SolutionTwo},
	{d04.SolutionOne, d04.SolutionTwo},
	{d05.Solution, d05.Solution},
	{d06.SolutionOne, d06.SolutionTwo},
	{d07.Solution, d07.Solution},
	{d08.SolutionOne, d08.SolutionTwo},
	{d09.SolutionOne, d09.SolutionTwo},
	{d10.Solution, d10.Solution},
	{d11.SolutionOne, d11.SolutionTwo},
	{d12.SolutionOne},
}

func main() {
	if err := run(); err != nil {
		println(err.Error())

		os.Exit(1)
	}
}

func run() error {
	var (
		day, solution int
		useSample     bool
	)

	set := flag.NewFlagSet("advent2024", flag.ContinueOnError)

	set.BoolVar(&useSample, "use-sample", false, "use the sample input")
	set.IntVar(&day, "day", 1, "day to run")
	set.IntVar(&solution, "solution", 1, "solution to run")

	args, _ := internal.Args(os.Args)

	err := set.Parse(args)
	if err != nil {
		return fmt.Errorf("invalid arguments: %w\n", err)
	}

	if day > len(days) {
		return fmt.Errorf("I didn't manage more than %d days this year...", len(days))
	}

	solutions := days[day-1]

	if solution > len(solutions) {
		return fmt.Errorf("I only have %d solutions for day %d", len(solutions), day)
	}

	dayDir := fmt.Sprintf("d%02d", day)

	var candidateFiles []string

	if useSample {
		candidateFiles = []string{
			fmt.Sprintf("sample-%d.txt", solution),
			"sample.txt",
		}
	} else {
		candidateFiles = []string{
			fmt.Sprintf("input-%d.txt", solution),
			"input.txt",
		}
	}

	var (
		input    io.Reader
		inputErr error
	)

	// Try each candidate in turn, joining the errors for one big soulfelt
	// error scream if we don't get what we want.
	for name := range slices.Values(candidateFiles) {
		inputPath := filepath.Join(dayDir, name)

		f, err := os.Open(inputPath)
		if err != nil {
			inputErr = errors.Join(inputErr, fmt.Errorf(
				"open input: %w", err))

			continue
		}

		input = f
		inputErr = nil

		break
	}

	if inputErr != nil {
		return fmt.Errorf("could not open input file: %w", inputErr)
	}

	fn := solutions[solution-1]

	err = fn(input)
	if err != nil {
		return fmt.Errorf("solution returned an error: %w", err)
	}

	return nil
}
