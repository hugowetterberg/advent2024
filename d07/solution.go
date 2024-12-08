package d07

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Solution(input io.Reader) error {
	sc := bufio.NewScanner(input)

	var (
		sumValidPart1 int
		sumValidPart2 int
		line          int
	)

	// Operators to use for part 1.
	operatorsPart1 := []operator{
		func(a, b int) int {
			return a * b
		},
		func(a, b int) int {
			return a + b
		},
	}

	// Operators to use for part 1. Here we add the strange concatenation
	// operator.
	operatorsPart2 := append([]operator{
		// Concatenate numbers
		func(a, b int) int {
			n, _ := strconv.Atoi(fmt.Sprintf("%d%d", a, b))
			return n
		},
	}, operatorsPart1...)

	for sc.Scan() {
		line++

		val, numStr, ok := strings.Cut(sc.Text(), ": ")
		if !ok {
			return fmt.Errorf("invalid input format on line %d", line)
		}

		value, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid value on line %d: %w",
				line, err)
		}

		ns := strings.Split(numStr, " ")
		numbers := make([]int, len(ns))

		for i := range ns {
			n, err := strconv.Atoi(ns[i])
			if err != nil {
				return fmt.Errorf(
					"invalid number at position %d on line %d: %w",
					i+1, line, err)
			}

			numbers[i] = n
		}

		if seek(numbers, value, operatorsPart1) {
			sumValidPart1 += value
		}

		if seek(numbers, value, operatorsPart2) {
			sumValidPart2 += value
		}
	}

	fmt.Printf(`Sum of valid equations:
  Part 1: %d
  Part 2: %d
`, sumValidPart1, sumValidPart2)

	return nil
}

type operator func(a, b int) int

func seek(numbers []int, target int, operators []operator) bool {
	return _seek(numbers, 0, 0, target, operators)
}

func _seek(numbers []int, idx int, prod int, target int, operators []operator) bool {
	// The product/sum/whatever always increases, so if we overshoot there's
	// no use in exploring the branch further.
	if prod > target {
		return false
	}

	// If we've hit the end of the numbers slice we return true if the
	// combination of operators resulted in the right product.
	if idx == len(numbers) {
		return prod == target
	}

	// If we're at the firs index there's no operation to be applied,
	// recursive call to index 1.
	if idx == 0 {
		return _seek(numbers, 1, numbers[0], target, operators)
	}

	// Iterate over the operators we have available.
	for _, op := range operators {
		// Make a recursive seek call for the next index, using the
		// result of `prod OP current_number` as the new product.
		found := _seek(numbers, idx+1, op(prod, numbers[idx]), target, operators)

		// If we found a branch that produces the correct result we
		// return true immediately, we're not interested in an
		// exhaustive search.
		if found {
			return true
		}
	}

	return false
}
