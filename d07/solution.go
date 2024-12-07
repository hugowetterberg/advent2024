package d07

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func SolutionOne(input io.Reader) error {
	sc := bufio.NewScanner(input)

	var (
		sumValid int
		line     int
	)

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

		if seek(numbers, value) {
			sumValid += value
		}
	}

	fmt.Printf("Sum of valid equations: %d\n", sumValid)

	return nil
}

type operator func(a, b int) int

var operators = []operator{
	func(a, b int) int {
		return a * b
	},
	func(a, b int) int {
		return a + b
	},
}

func seek(numbers []int, target int) bool {
	return _seek(numbers, 0, 0, target)
}

func _seek(numbers []int, idx int, prod int, target int) bool {
	if prod > target {
		return false
	}

	if idx == len(numbers) {
		return prod == target
	}

	if idx == 0 {
		return _seek(numbers, 1, numbers[0], target)
	}

	for _, op := range operators {
		found := _seek(numbers, idx+1, op(prod, numbers[idx]), target)
		if found {
			return true
		}
	}

	return false
}
