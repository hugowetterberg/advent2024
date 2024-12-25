package d12

import (
	"bufio"
	"fmt"
	"io"
	"slices"
)

type Vec struct {
	Row int
	Col int
}

func (v Vec) Add(b Vec) Vec {
	return Vec{
		Row: v.Row + b.Row,
		Col: v.Col + b.Col,
	}
}

func (v Vec) Subtract(b Vec) Vec {
	return Vec{
		Row: v.Row - b.Row,
		Col: v.Col - b.Col,
	}
}

func (v Vec) String() string {
	return fmt.Sprintf("{r: %d, c: %d}", v.Row, v.Col)
}

// Uppercase ASCII doesn't use the last bit, use that to store a flag.
const VisitedFlag = 1 << 7

type Region struct {
	Area      int
	Perimeter int
}

func MeasureRegion(start Vec, garden [][]byte) Region {
	var r Region

	rPlant, visited, inside := getPlot(garden, start)
	if visited || !inside {
		return Region{}
	}

	r.Area, r.Perimeter = _visit(garden, rPlant, start)

	return r
}

func _visit(garden [][]byte, rPlant byte, pos Vec) (int, int) {
	var area, perimeter, nCount int

	garden[pos.Row][pos.Col] |= VisitedFlag

	for i := range neighbours {
		np := pos.Add(neighbours[i])

		plant, visited, inside := getPlot(garden, np)
		if !inside || plant != rPlant {
			continue
		}

		nCount++

		if visited {
			continue
		}

		a, p := _visit(garden, rPlant, np)

		area += a
		perimeter += p
	}

	return area + 1, perimeter + (4 - nCount)
}

// getPlot returns the plant type, visited state and whether the position is
// within the garden.
func getPlot(garden [][]byte, pos Vec) (byte, bool, bool) {
	if pos.Row < 0 || pos.Col < 0 {
		return 0, false, false
	}

	if pos.Row >= len(garden) || pos.Col >= len(garden[pos.Row]) {
		return 0, false, false
	}

	v := garden[pos.Row][pos.Col]

	return v &^ VisitedFlag, v&VisitedFlag == VisitedFlag, true
}

var neighbours = []Vec{
	{Row: -1},
	{Col: 1},
	{Row: 1},
	{Col: -1},
}

func SolutionOne(input io.Reader) error {
	var garden [][]byte

	sc := bufio.NewScanner(input)

	for sc.Scan() {
		garden = append(garden, slices.Clone(sc.Bytes()))
	}

	err := sc.Err()
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	var sum int

	for r := range garden {
		for c := range garden[r] {
			r := MeasureRegion(Vec{
				Row: r,
				Col: c,
			}, garden)

			sum += r.Area * r.Perimeter
		}
	}

	fmt.Printf("Total price of fencing: %d\n", sum)

	return nil
}
