package d08

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
	"strings"
)

var towerChars = "0123456789abcdefghijklmnopqrstuvwxuzABCDEFGHIJKLMNOPQRSTUVWXUZ"

func towerFreq(b byte) uint64 {
	idx := strings.Index(towerChars, string(b))
	if idx == -1 {
		return 0
	}

	return 1 << idx
}

func towerChar(freq uint64) byte {
	for i := range uint64(len(towerChars)) {
		flag := uint64(1) << i

		if freq&flag == flag {
			return towerChars[i]
		}
	}

	return 0
}

func SolutionOne(input io.Reader) error {
	return solution(input, false)
}

func SolutionTwo(input io.Reader) error {
	return solution(input, true)
}

func solution(input io.Reader, resonantHarmonics bool) error {
	var antinodes FlagMap
	var towerMap FlagMap

	var antennas []Antenna

	sc := bufio.NewScanner(input)

	for sc.Scan() {
		data := sc.Bytes()

		tLine := make([]uint64, len(data))
		aLine := make([]uint64, len(data))

		for i, b := range data {
			code := towerFreq(b)
			if code == 0 {
				continue
			}

			tLine[i] = code

			antennas = append(antennas, Antenna{
				Freq: code,
				Position: Coord{
					Row: len(antinodes),
					Col: i,
				},
			})
		}

		towerMap = append(towerMap, tLine)
		antinodes = append(antinodes, aLine)
	}

	for i, a := range antennas {
		for j, b := range antennas {
			if i == j || a.Freq != b.Freq {
				continue
			}

			if resonantHarmonics {
				antinodes.Set(a.Position, a.Freq)
			}

			delta := b.Position.Subtract(a.Position)

			m := 2

			for {
				aPos := a.Position.Add(delta.Mult(m))

				if !antinodes.Set(aPos, a.Freq) || !resonantHarmonics {
					break
				}

				m++
			}

		}
	}

	var antinodeCount int

	for i := range antinodes {
		for j := range antinodes[i] {
			if antinodes[i][j] != 0 {
				antinodeCount++
			}
		}
	}

	DebugDump(towerMap, antinodes)

	fmt.Printf("Number of antinodes: %d\n", antinodeCount)

	return nil
}

func expectCount(n int) int {
	return n * (n - 1)
}

func DebugDump(towers FlagMap, antinodes FlagMap) {
	for row := range towers {
		for col := range towers[row] {
			switch {
			case antinodes[row][col] != 0:
				print("#")
			case towers[row][col] != 0:
				print(string(towerChar(towers[row][col])))
			default:
				print(".")
			}
		}

		println()
	}
}

func DebugDumpMask(towers FlagMap, antinodes FlagMap, freq uint64) {
	for row := range towers {
		for col := range towers[row] {
			switch {
			case antinodes[row][col]&freq == freq:
				print("#")
			case towers[row][col]&freq == freq:
				print(string(towerChar(towers[row][col])))
			default:
				print(".")
			}
		}

		println()
	}
}

type Coord struct {
	Row int
	Col int
}

func (c Coord) Mult(v int) Coord {
	return Coord{
		Row: c.Row * v,
		Col: c.Col * v,
	}
}

func (c Coord) Add(v Coord) Coord {
	return Coord{
		Row: c.Row + v.Row,
		Col: c.Col + v.Col,
	}
}

func (c Coord) Subtract(v Coord) Coord {
	return Coord{
		Row: c.Row - v.Row,
		Col: c.Col - v.Col,
	}
}

type Antenna struct {
	Freq     uint64
	Position Coord
}

type FlagMap [][]uint64

func outsideBounds(fm FlagMap, pos Coord) bool {
	return pos.Row < 0 || pos.Col < 0 ||
		pos.Row >= len(fm) || pos.Col >= len(fm[pos.Row])
}

func (fm FlagMap) CountAllMask(mask uint64) int {
	var count int

	for i := range fm {
		for j := range fm[i] {
			if fm[i][j]&mask == mask {
				count++
			}
		}
	}

	return count
}

func (fm FlagMap) Get(pos Coord) (uint64, bool) {
	if outsideBounds(fm, pos) {
		return 0, false
	}

	return fm[pos.Row][pos.Col], true
}

func (fm FlagMap) Set(pos Coord, code uint64) bool {
	if outsideBounds(fm, pos) {
		return false
	}

	fm[pos.Row][pos.Col] |= code

	return true
}

func (fm FlagMap) Count(pos Coord) int {
	if outsideBounds(fm, pos) {
		return 0
	}

	return bits.OnesCount64(fm[pos.Row][pos.Col])
}
