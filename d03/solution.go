package d03

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var mulExp = regexp.MustCompile(`mul\((\d+),(\d+)\)`)

func SolutionOne(input io.Reader) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	matches := mulExp.FindAllStringSubmatch(string(data), -1)

	var sum int

	for _, sub := range matches {
		// Ignoring errors, we already know from the regex that we're
		// dealing with digits only.
		a, _ := strconv.Atoi(sub[1])
		b, _ := strconv.Atoi(sub[2])

		sum += a * b
	}

	fmt.Printf("Sum: %d\n", sum)

	return nil
}

var mulCondExp = regexp.MustCompile(`(mul\((\d+),(\d+)\))|(do\(\)|don't\(\))`)

func SolutionTwo(input io.Reader) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	matches := mulCondExp.FindAllStringSubmatch(string(data), -1)

	var (
		sum      int
		disabled bool
	)

	for _, sub := range matches {
		if strings.HasPrefix(sub[0], "do") {
			disabled = sub[0] == "don't()"
		}

		if disabled {
			continue
		}

		// Ignoring errors, we already know from the regex that we're
		// dealing with digits only.
		a, _ := strconv.Atoi(sub[2])
		b, _ := strconv.Atoi(sub[3])

		sum += a * b
	}

	fmt.Printf("Sum: %d\n", sum)

	return nil
}
