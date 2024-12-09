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
