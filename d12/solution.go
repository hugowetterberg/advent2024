package d12

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io"
	"slices"

	"github.com/hugowetterberg/advent2024/internal"
	"github.com/llgcode/draw2d/draw2dimg"
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

func (v Vec) MultInt(m int) Vec {
	return Vec{
		Row: v.Row * m,
		Col: v.Col * m,
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

type Line struct {
	drop bool

	Start    Vec
	Length   int
	Norm     Vec
	Position Vec
}

func (l Line) String() string {
	return fmt.Sprintf("%s->%s", l.Start, l.End())
}

func (l Line) End() Vec {
	return l.Start.Add(l.Norm.MultInt(l.Length))
}

func (l Line) TryJoin(b Line) (Line, bool) {
	if l.Norm != b.Norm || l.Position != b.Position {
		return Line{}, false
	}

	switch {
	case l.End() == b.Start:
		l.Length += b.Length

		return l, true
	case b.End() == l.Start:
		b.Length += l.Length

		return b, true
	default:
		return Line{}, false
	}
}

func BorderLine(pos, neighbour Vec) Line {
	dir := neighbour.Subtract(pos)

	offset := borderOffset[dir]

	l := Line{
		Start:    pos.Add(offset),
		Length:   1,
		Position: dir,
	}

	switch dir {
	case up, down:
		l.Norm = Vec{0, 1}
	case left, right:
		l.Norm = Vec{1, 0}
	}

	return l
}

// Uppercase ASCII doesn't use the last bit, use that to store a flag.
const VisitedFlag = 1 << 7

type Region struct {
	Plant     byte
	Area      int
	Perimeter int
	Stretches []Line
}

func (r *Region) AddStretch(s Line) {
	for i, l := range r.Stretches {
		n, ok := l.TryJoin(s)
		if !ok {
			continue
		}

		r.Stretches[i] = n

		return
	}

	r.Stretches = append(r.Stretches, s)
}

func (r *Region) Simplify() {
	var dropSum int

	for {
		var dropCount int

		for i, l := range r.Stretches {
			if l.drop {
				continue
			}

			for j, b := range r.Stretches {
				if j == i || b.drop {
					continue
				}

				n, ok := l.TryJoin(b)
				if !ok {
					continue
				}

				r.Stretches[j].drop = true
				r.Stretches[i] = n

				dropCount++
				dropSum++
			}
		}

		if dropCount == 0 {
			break
		}
	}

	stretches := make([]Line, 0, len(r.Stretches)-dropSum)

	for _, l := range r.Stretches {
		if l.drop {
			continue
		}

		stretches = append(stretches, l)
	}

	r.Stretches = stretches
}

func MeasureRegion(start Vec, garden [][]byte) (Region, bool) {
	rPlant, visited, inside := getPlot(garden, start)
	if visited || !inside {
		return Region{}, false
	}

	r := Region{
		Plant: rPlant,
	}

	r.Area, r.Perimeter = _visit(garden, &r, start)

	r.Simplify()

	return r, true
}

func _visit(garden [][]byte, r *Region, pos Vec) (int, int) {
	var area, perimeter, nCount int

	garden[pos.Row][pos.Col] |= VisitedFlag

	for i := range neighbours {
		np := pos.Add(neighbours[i])

		plant, visited, inside := getPlot(garden, np)
		if !inside || plant != r.Plant {
			// Border baby!
			r.AddStretch(BorderLine(pos, np))

			continue
		}

		nCount++

		if visited {
			continue
		}

		a, p := _visit(garden, r, np)

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

var (
	up    = Vec{Row: -1, Col: 0}
	right = Vec{Row: 0, Col: 1}
	down  = Vec{Row: 1, Col: 0}
	left  = Vec{Row: 0, Col: -1}
)

var neighbours = []Vec{
	up,
	right,
	down,
	left,
}

var borderOffset = map[Vec]Vec{
	up:    {Row: 0, Col: 0},
	right: {Row: 0, Col: 1},
	down:  {Row: 1, Col: 0},
	left:  {Row: 0, Col: 0},
}

const debugGrid = 20

func Solution(input io.Reader) error {
	set := flag.NewFlagSet("day12", flag.ContinueOnError)

	var (
		debugRegion string
		verbose     bool
	)

	set.BoolVar(&verbose, "verbose", false, "verbose region output")
	set.StringVar(&debugRegion, "debug-region", "", "output debug image for region")

	err := internal.ParseSolutionFlags(set)
	if err != nil {
		return err
	}

	var garden [][]byte

	sc := bufio.NewScanner(input)

	for sc.Scan() {
		garden = append(garden, slices.Clone(sc.Bytes()))
	}

	err = sc.Err()
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	width := len(garden[0])
	height := len(garden)

	var (
		dest *image.RGBA
		gc   *draw2dimg.GraphicContext
	)

	if debugRegion != "" {
		dest = image.NewRGBA(image.Rect(0, 0,
			width*debugGrid+debugGrid*2,
			height*debugGrid+debugGrid*2))
		gc = draw2dimg.NewGraphicContext(dest)

		gc.SetFillColor(color.Black)
		gc.Clear()
	}

	var sum, bulkPrice int

	for r := range garden {
		for c := range garden[r] {
			r, ok := MeasureRegion(Vec{
				Row: r,
				Col: c,
			}, garden)
			if !ok {
				continue
			}

			sum += r.Area * r.Perimeter
			bulkPrice += r.Area * len(r.Stretches)

			if verbose {
				fmt.Printf("%s Area: %d Perimeter: %d\n",
					[]byte{r.Plant}, r.Area, r.Perimeter)

				for _, s := range r.Stretches {
					println(s.String())
				}

				println()
			}

			if gc != nil && debugRegion[0] == r.Plant {
				drawRegion(gc, r)
			}
		}
	}

	fmt.Printf("Standard price of fencing: %d\n", sum)
	fmt.Printf("Bulk price of fencing: %d\n", bulkPrice)

	if gc != nil {
		draw2dimg.SaveToPngFile("debug.png", dest)
	}

	return nil
}

func drawRegion(gc *draw2dimg.GraphicContext, r Region) {
	for i, s := range r.Stretches {
		gc.SetStrokeColor(color.RGBA{
			uint8(100 + (i*20)%156),
			uint8(100 + (+i*40)%156),
			uint8(100 + (i*80)%156),
			0xff,
		})

		switch s.Position {
		case up, left:
			gc.SetLineWidth(5)
		case down, right:
			gc.SetLineWidth(1)
		}

		gc.BeginPath()

		start := vec2coord(s.Start)
		end := vec2coord(s.Start.Add(s.Norm.MultInt(s.Length)))

		gc.MoveTo(start[0], start[1])
		gc.LineTo(end[0], end[1])

		gc.Close()
		gc.FillStroke()
	}
}

func vec2coord(v Vec) [2]float64 {
	return [2]float64{
		float64(debugGrid + v.Col*debugGrid),
		float64(debugGrid + v.Row*debugGrid),
	}
}
