package d11

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// If the stone is engraved with the number 0, it is replaced by a stone
// engraved with the number 1.
//
// If the stone is engraved with a number that has an even number of digits, it
// is replaced by two stones. The left half of the digits are engraved on the
// new left stone, and the right half of the digits are engraved on the new
// right stone. (The new numbers don't keep extra leading zeroes: 1000 would
// become stones 10 and 0.)
//
// If none of the other rules apply, the stone is replaced by a new stone; the
// old stone's number multiplied by 2024 is engraved on the new stone.

type Stone uint64

var stoneExp = []Stone{
	1e1, 1e2, 1e3, 1e4, 1e5,
	1e6, 1e7, 1e8, 1e9, 1e10,
	1e11, 1e12, 1e13, 1e14, 1e15,
	1e16, 1e17, 1e18, 1e19,
}

func (s Stone) Digits() int {
	for i := 0; i < len(stoneExp); i++ {
		if s < stoneExp[i] {
			return i + 1
		}
	}

	return 20
}

func (s Stone) Split() (Stone, Stone, bool) {
	d := s.Digits()
	if d%2 != 0 {
		return 0, 0, false
	}

	// 1234 = 4 digits
	// 1234 / 1e2 = 12
	// 1234 % 1e2 = 34

	splitExp := stoneExp[d/2-1]

	return s / splitExp, s % splitExp, true
}

func SolutionOne(input io.Reader) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	stoneStrs := strings.Split(
		string(bytes.TrimSpace(data)), " ")

	stones := make([]Stone, len(stoneStrs))

	for i := range stoneStrs {
		s, err := strconv.ParseUint(stoneStrs[i], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid stone at position %d: %w",
				i+1, err)
		}

		stones[i] = Stone(s)
	}

	var back, front []Stone

	front = stones

	for range 25 {
		for _, s := range front {
			if s == 0 {
				back = append(back, 1)
			} else if h, l, ok := s.Split(); ok {
				back = append(back, h, l)
			} else {
				back = append(back, s*2024)
			}
		}

		front, back = back, front
		back = back[0:0]
	}

	fmt.Printf("Number of stones: %v\n", len(front))

	return nil
}