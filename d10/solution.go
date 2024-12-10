package d10

import (
	"bufio"
	"fmt"
	"io"
	"slices"
)

type TopoMap [][]byte

func outsideBounds(m TopoMap, pos Coord) bool {
	return pos[0] < 0 || pos[1] < 0 ||
		pos[0] >= len(m) || pos[1] >= len(m[pos[0]])
}

func (m TopoMap) GetHeight(pos Coord) (byte, bool) {
	if outsideBounds(m, pos) {
		return 0, false
	}

	return m[pos[1]][pos[0]], true
}

type Coord [2]int

func (c Coord) Add(v Coord) Coord {
	return Coord{
		c[0] + v[0],
		c[1] + v[1],
	}
}

var directions = []Coord{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

func SolutionOne(input io.Reader) error {
	var topo TopoMap

	sc := bufio.NewScanner(input)

	for sc.Scan() {
		row := slices.Clone(sc.Bytes())

		for i, b := range row {
			row[i] = b - 48
		}

		topo = append(topo, row)
	}

	err := sc.Err()
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	var sum int

	for y := range topo {
		for x := range topo[y] {
			if topo[y][x] != 0 {
				continue
			}

			found := explore(topo, Coord{x, y})

			score := len(found)

			sum += score
		}
	}

	fmt.Printf("Sum of trailhead scores: %d\n", sum)

	return nil
}

func explore(m TopoMap, pos Coord) []Coord {
	h, ok := m.GetHeight(pos)
	if !ok {
		return nil
	}

	return _explore(m, pos, h)
}

func _explore(m TopoMap, pos Coord, h byte) []Coord {
	if h == 9 {
		return []Coord{pos}
	}

	var endpoints []Coord

	for i := range directions {
		newPos := pos.Add(directions[i])

		newHeight, ok := m.GetHeight(newPos)
		if !ok || newHeight != h+1 {
			continue
		}

		discovered := _explore(m, newPos, newHeight)

		for _, p := range discovered {
			if slices.Contains(endpoints, p) {
				continue
			}

			endpoints = append(endpoints, p)
		}
	}

	return endpoints
}
