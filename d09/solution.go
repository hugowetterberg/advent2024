package d09

import (
	"bytes"
	"fmt"
	"io"
	"slices"
)

func SolutionOne(input io.Reader) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	var disk []int

	for i, b := range bytes.TrimSpace(data) {
		if b < 48 || b > 57 {
			return fmt.Errorf("invalid character %q at position %d",
				string(b), i+1)
		}

		n := int(b - 48)
		isFile := i%2 == 0

		if isFile {
			id := i / 2

			disk = append(disk, slices.Repeat([]int{id}, n)...)
		} else {
			disk = append(disk, slices.Repeat([]int{-1}, n)...)
		}
	}

	var head int

	for tail := len(disk) - 1; tail > 0; tail-- {
		for disk[head] != -1 && head < tail {
			head++
		}

		for disk[tail] == -1 && tail > head {
			tail--
		}

		if tail == head {
			break
		}

		disk[head], disk[tail] = disk[tail], disk[head]
	}

	var checksum uint64

	for i, n := range disk {
		if n == -1 {
			break
		}

		checksum += uint64(i * n)
	}

	fmt.Printf("Checksum is: %d\n", checksum)

	return nil
}

func SolutionTwo(input io.Reader) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	var disk []int
	var free []DiskRef

	for i, b := range bytes.TrimSpace(data) {
		if b < 48 || b > 57 {
			return fmt.Errorf("invalid character %q at position %d",
				string(b), i+1)
		}

		n := int(b - 48)
		isFile := i%2 == 0

		if isFile {
			id := i / 2

			disk = append(disk, slices.Repeat([]int{id}, n)...)
		} else {
			free = append(free, DiskRef{
				len(disk), n,
			})

			disk = append(disk, slices.Repeat([]int{-1}, n)...)
		}
	}

	for tail := len(disk) - 1; tail > 0; tail-- {
		for disk[tail] == -1 && tail > 0 {
			tail--
		}

		size := 1
		id := disk[tail]

		// Find beginning of file.
		for tail > 0 && disk[tail-1] == id {
			size++
			tail--
		}

		ref := -1

		for r := 0; r < len(free); r++ {
			if free[r][1] >= size {
				ref = r
				break
			}
		}

		if ref < 0 || free[ref][0] > tail {
			continue
		}

		copy(disk[free[ref][0]:], disk[tail:tail+size])

		for i := 0; i < size; i++ {
			disk[tail+i] = -1
		}

		free[ref][0] += size
		free[ref][1] -= size
	}

	var checksum uint64

	for i, n := range disk {
		if n == -1 {
			continue
		}

		checksum += uint64(i * n)
	}

	fmt.Printf("Checksum is: %d\n", checksum)

	return nil
}

type DiskRef [2]int
