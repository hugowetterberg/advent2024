package d06

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type MapState byte

const (
	MapStateNone     = 0
	MapStateVisited  = 1
	MapStateBlocked  = 2
	MapStateDirStart = 4
)

type AreaMap [][]MapState

func (am AreaMap) inBounds(pos Vec) bool {
	return pos[0] >= 0 && pos[1] >= 0 &&
		pos[1] < len(am) && pos[0] < len(am[pos[1]])
}

func (am AreaMap) Check(pos Vec) (MapState, bool) {
	if !am.inBounds(pos) {
		return MapStateNone, false
	}

	return am[pos[1]][pos[0]], true
}

func (am AreaMap) Get(pos Vec) MapState {
	if !am.inBounds(pos) {
		return MapStateNone
	}

	return am[pos[1]][pos[0]]
}

func (am AreaMap) Set(pos Vec, state MapState) {
	if !am.inBounds(pos) {
		return
	}

	am[pos[1]][pos[0]] |= state
}

func (am AreaMap) Unset(pos Vec, state MapState) {
	if !am.inBounds(pos) {
		return
	}

	cs := am[pos[1]][pos[0]]
	if cs&state == state {
		am[pos[1]][pos[0]] -= state
	}
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
			case cell&MapStateVisited == MapStateVisited:
				print("X")
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
	state, err := ReadState(input)
	if err != nil {
		return err
	}

	for state.Move() {
	}

	visitedCount := state.Area.Count(MapStateVisited)

	fmt.Printf("Visited %d locations\n", visitedCount)

	return nil
}

func SolutionTwo(input io.Reader) error {
	state, err := ReadState(input)
	if err != nil {
		return err
	}

	var loopCount int

	for y := range state.Area {
		for x, cell := range state.Area[y] {
			if cell&MapStateBlocked == MapStateBlocked {
				continue
			}

			isLoop := blockAndCheckForLoop(*state, Vec{x, y})
			if isLoop {
				loopCount++
			}
		}
	}

	fmt.Printf("Loop location count: %d\n", loopCount)

	return nil
}

func blockAndCheckForLoop(s State, pos Vec) bool {
	s.Area.Set(pos, MapStateBlocked)
	defer func() {
		s.Area.Unset(pos, MapStateBlocked)

		// Reset to only blocked states.
		for y := range s.Area {
			for x, cell := range s.Area[y] {
				s.Area[y][x] = cell & MapStateBlocked
			}
		}
	}()

	for s.Move() {
		current := s.Area.Get(s.GuardPos)

		// Dynamic direction state that we use to mark traversal in a
		// given direction.
		dirState := MapState(MapStateDirStart << s.GuardDir)

		// Check if we've visited this position in the current direction.
		if current&dirState == dirState {
			return true
		}

		s.Area.Set(s.GuardPos, dirState)
	}

	return false
}

type State struct {
	Area     AreaMap
	GuardPos Vec
	GuardDir int
}

func (s *State) Move() bool {
	nextPos := s.GuardPos.Add(directions[s.GuardDir])
	nextState, insideArea := s.Area.Check(nextPos)

	switch {
	case !insideArea:
		s.GuardPos = nextPos
		return false
	case nextState&MapStateBlocked == MapStateBlocked:
		s.GuardDir = (s.GuardDir + 1) % len(directions)
	default:
		s.GuardPos = nextPos
		s.Area.Set(s.GuardPos, MapStateVisited)
	}

	return true
}

func (s *State) DebugDump(w io.Writer) {
	s.Area.DebugDump(w, s.GuardPos, s.GuardDir)
}

func ReadState(input io.Reader) (*State, error) {
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
					return nil, fmt.Errorf(
						"unknown guard direction %q",
						string(line[i]))
				}

				guardPos = Vec{i, len(area)}
			}
		}

		area = append(area, mapLine)
	}

	area.Set(guardPos, MapStateVisited)

	return &State{
		Area:     area,
		GuardPos: guardPos,
		GuardDir: dirIdx,
	}, nil
}
