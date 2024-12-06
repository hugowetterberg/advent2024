package d06

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type MapState byte

const (
	MapStateNone    = 0
	MapStateVisited = 1
	MapStateBlocked = 2
)

type AreaMap [][]MapState

func (am AreaMap) Check(pos Vec) (MapState, bool) {
	if pos[1] >= len(am) || pos[0] >= len(am[pos[1]]) {
		return MapStateNone, false
	}

	return am[pos[1]][pos[0]], true
}

func (am AreaMap) Set(pos Vec, state MapState) {
	if pos[1] >= len(am) || pos[0] >= len(am[pos[1]]) {
		return
	}

	am[pos[1]][pos[0]] |= state
}

func (am AreaMap) Count(state MapState) int {
	var count int

	for i := range am {
		for j := range am[i] {
			if am[i][j]&state == state {
				count++
			}
		}
	}

	return count
}

func (am AreaMap) DebugDump(w io.Writer, guardPos Vec, dirIdx int) {
	for y, row := range am {
		for x, cell := range row {
			switch {
			case guardPos[0] == x && guardPos[1] == y:
				print(string(dirSymbols[dirIdx]))
			case cell&MapStateBlocked == MapStateBlocked:
				print("#")
			default:
				print(".")
			}
		}

		println()
	}
}

type Vec [2]int

func (v Vec) Add(v2 Vec) Vec {
	return Vec{v[0] + v2[0], v[1] + v2[1]}
}

var dirSymbols = "^>v<"

var directions = []Vec{
	{0, -1}, {1, 0}, {0, 1}, {-1, 0},
}

func SolutionOne(input io.Reader) error {
	sc := bufio.NewScanner(input)

	var (
		area     AreaMap
		dirIdx   int
		guardPos Vec
	)

	for sc.Scan() {
		line := sc.Bytes()

		mapLine := make([]MapState, len(line))

		for i := range line {
			switch line[i] {
			case '.':
			case '#':
				mapLine[i] = MapStateBlocked
			default:
				dirIdx = strings.IndexByte(dirSymbols, line[i])
				if dirIdx == -1 {
					return fmt.Errorf("unknown guard direction %q", string(line[i]))
				}

				guardPos = Vec{i, len(area)}
			}
		}

		area = append(area, mapLine)
	}

	area.Set(guardPos, MapStateVisited)

	for {
		nextPos := guardPos.Add(directions[dirIdx])

		nextState, insideArea := area.Check(nextPos)
		if !insideArea {
			guardPos = nextPos
			break
		}

		if nextState&MapStateBlocked == MapStateBlocked {
			dirIdx = (dirIdx + 1) % len(directions)
		} else {
			guardPos = nextPos
			area.Set(guardPos, MapStateVisited)
		}
	}

	visitedCount := area.Count(MapStateVisited)

	fmt.Printf("Visited %d locations\n", visitedCount)

	return nil
}
