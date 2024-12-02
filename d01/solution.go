package d01

import (
	"errors"
	"fmt"
	"io"
	"slices"
)

func SolutionOne(input io.Reader) error {
	listA, listB, err := readLists(input)
	if err != nil {
		return fmt.Errorf("read lists: %w", err)
	}

	if len(listA) != len(listB) {
		return errors.New("expected lists to have the same lenght")
	}

	slices.Sort(listA)
	slices.Sort(listB)

	var distanceSum int

	for i := range listA {
		distanceSum += max(listA[i], listB[i]) - min(listA[i], listB[i])
	}

	fmt.Printf("Sum: %d\n", distanceSum)

	return nil
}

func SolutionTwo(input io.Reader) error {
	listA, listB, err := readLists(input)
	if err != nil {
		return fmt.Errorf("read lists: %w", err)
	}

	var similarity int

	for vA := range slices.Values(listA) {
		var count int

		for vB := range slices.Values(listB) {
			if vB == vA {
				count++
			}
		}

		similarity += vA * count
	}

	fmt.Printf("Similarity: %d\n", similarity)

	return nil
}

func readLists(r io.Reader) ([]int, []int, error) {
	var (
		a, b   []int
		aN, bN int
		line   int
	)

	for {
		line++

		_, err := fmt.Fscanf(r, "%d   %d\n", &aN, &bN)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to read line %d: %w",
				line, err)
		}

		a = append(a, aN)
		b = append(b, bN)
	}

	return a, b, nil
}
