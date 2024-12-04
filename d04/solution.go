package d04

import (
	"bufio"
	"fmt"
	"io"
)

type Vec [2]int

func (v Vec) Add(a Vec) Vec {
	return Vec{
		v[0] + a[0],
		v[1] + a[1],
	}
}

type Lines [][]byte

func (l Lines) CharAt(pos Vec) (byte, bool) {
	// Bounds check
	if pos[0] < 0 ||
		pos[1] < 0 ||
		pos[0] >= len(l) ||
		pos[1] >= len(l[pos[0]]) {
		return 0, false
	}

	return l[pos[0]][pos[1]], true
}

func SolutionOne(input io.Reader) error {
	directions := []Vec{
		{0, -1},
		{1, -1},
		{1, 0},
		{1, 1},
		{0, 1},
		{-1, 1},
		{-1, 0},
		{-1, -1},
	}

	var lines Lines

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := make([]byte, len(scanner.Bytes()))

		copy(line, scanner.Bytes())

		lines = append(lines, line)
	}

	err := scanner.Err()
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	word := []byte("XMAS")

	var count int

	for row := range len(lines) {
		for col := range len(lines[row]) {
			for _, dir := range directions {
				hit := checkWord(
					lines,
					Vec{row, col},
					dir,
					word)
				if hit {
					count++
				}
			}
		}
	}

	fmt.Printf("Words found: %d\n", count)

	return nil
}

func checkWord(lines Lines, pos Vec, direction Vec, word []byte) bool {
	if len(word) == 0 {
		return true
	}

	char, ok := lines.CharAt(pos)
	if !ok {
		return false
	}

	if word[0] != char {
		return false
	}

	return checkWord(lines, pos.Add(direction), direction, word[1:])
}
