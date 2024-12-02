package d02

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func SolutionOne(input io.Reader) error {
	lines := bufio.NewScanner(input)

	var linum, safeCount int

	for lines.Scan() {
		linum++

		numberStrings := strings.Split(lines.Text(), " ")

		var (
			current    int
			currentDir int
			unsafe     bool
		)

		for i, nStr := range numberStrings {
			n, err := strconv.Atoi(nStr)
			if err != nil {
				return fmt.Errorf(
					"parse number %d on line %d: %w",
					i+1, linum, err)
			}

			if i == 0 {
				current = n
				continue
			}

			diff := delta(n, current)
			if diff == 0 || diff > 3 {
				unsafe = true

				break
			}

			dir := n - current
			if (dir < 0 && currentDir > 0) || (dir > 0 && currentDir < 0) {
				unsafe = true

				break
			}

			current = n
			currentDir = dir
		}

		if !unsafe {
			safeCount++
		}
	}

	err := lines.Err()
	if err != nil {
		return fmt.Errorf("read input lines: %w", err)
	}

	fmt.Printf("Safe reports: %d\n", safeCount)

	return nil
}

func SolutionTwo(input io.Reader) error {
	lines := bufio.NewScanner(input)

	var linum, safeCount int

	for lines.Scan() {
		linum++

		numberStrings := strings.Split(lines.Text(), " ")
		report := make([]int, len(numberStrings))

		for i, nStr := range numberStrings {
			n, err := strconv.Atoi(nStr)
			if err != nil {
				return fmt.Errorf(
					"parse number %d on line %d: %w",
					i+1, linum, err)
			}

			report[i] = n
		}

		unsafeIdx := findUnsafe(report)
		if unsafeIdx == -1 {
			safeCount++
			continue
		}

		alternatives := make([][]int, len(report))

		// Tried to just create two variants based on just removing the
		// offending value or the preceding value, but that heuristic
		// didn't work, so brute force it is.
		for i := range report {
			alternatives[i] = copyExcept(report, i)
		}

		if anySafe(alternatives...) {
			safeCount++
			continue
		}
	}

	err := lines.Err()
	if err != nil {
		return fmt.Errorf("read input lines: %w", err)
	}

	fmt.Printf("Safe reports: %d\n", safeCount)

	return nil
}

func copyExcept(src []int, idx int) []int {
	dst := make([]int, len(src)-1)

	if idx != 0 {
		copy(dst, src[:idx])
	}

	copy(dst[idx:], src[idx+1:])

	return dst
}

func anySafe(reports ...[]int) bool {
	for i := range reports {
		if findUnsafe(reports[i]) == -1 {
			return true
		}
	}

	return false
}

func findUnsafe(report []int) int {
	var (
		current    int
		currentDir int
	)

	for i, n := range report {
		if i == 0 {
			current = n
			continue
		}

		diff := delta(n, current)
		if diff == 0 || diff > 3 {
			return i
		}

		dir := n - current
		if (dir < 0 && currentDir > 0) || (dir > 0 && currentDir < 0) {
			return i
		}

		current = n
		currentDir = dir
	}

	return -1
}

func delta(a, b int) int {
	return abs(a - b)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}
