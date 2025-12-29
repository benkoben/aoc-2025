package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type grid struct {
	vertex [][]int
	edges  [][]int
	layout [][]string

	compressedX map[int]int
	compressedY map[int]int

	compressedXReversed map[int]int
	compressedYReversed map[int]int
}

func (g grid) String() string {
	out := ""

	for _, row := range g.layout {
		out = out + fmt.Sprintf("%v\n", row)
	}
	return out
}

func newGrid(input [][]int) grid {
	var g grid
	var rowCount int
	var columnCount int

	g.compressedX, g.compressedY, rowCount, columnCount = compress(input)
	compressed := [][]int{}
	for _, item := range input {
		x, y := item[0], item[1]
		compressed = append(compressed, []int{
			g.compressedX[x],
			g.compressedY[y],
		})
	}

	g.compressedXReversed = map[int]int{}
	g.compressedYReversed = map[int]int{}

	for k, v := range g.compressedX {
		g.compressedXReversed[v] = k
	}

	for k, v := range g.compressedY {
		g.compressedYReversed[v] = k
	}

	fmt.Println(rowCount, "x", columnCount)

	layout := make([][]string, rowCount)
	for y := range rowCount {
		row := make([]string, columnCount)
		layout[y] = row
	}

	for _, coord := range compressed {
		g.vertex = append(g.vertex, coord)
		x, y := coord[0], coord[1]
		layout[y][x] = "#"
	}

	edges := rasterize(compressed)
	for _, coord := range edges {
		g.edges = append(g.edges, coord)
		x, y := coord[0], coord[1]
		layout[y][x] = "#"
	}

	g.layout = layout
	g.fill()
	return g
}

func (g *grid) lookupOriginalCoordinate(x, y int) []int {
	var out []int

	out = append(out, g.compressedXReversed[x], g.compressedYReversed[y])

	return out
}

func (g *grid) insidePolygon(coords [][]int) bool {

	var out bool = true

	for _, coord := range coords {
		x, y := coord[0], coord[1]
		if g.layout[y][x] != "#" {
			out = false
		}
	}
	return out
}

func (g *grid) fill() {
	// raycast to find a start position
	start := []int{0, 0}

rayCast:
	for y := range g.layout {
		start[1] = y
		for x := range g.layout[y] {
			start[0] = x
			current := g.layout[y][x]
			if current == "" {
				lookAhead := g.layout[y][x:]
				if countOccurence("#", lookAhead) == 1 {
					break rayCast
				}
			}
		}
	}
	fmt.Println("start", start)
	g.dfs(start)
}

func (g *grid) dfs(s []int) {
	x, y := s[0], s[1]

	if g.layout[y][x] == "#" {
		fmt.Println()
		return
	}
	g.layout[y][x] = "#"

	adj := g.findAdj(s)
	// fmt.Println(adj)
	for _, coord := range adj {
		g.dfs(coord)
	}
}

func (g *grid) findAdj(s []int) [][]int {
	var adj [][]int

	maxX := len(g.layout[0]) - 1
	maxY := len(g.layout) - 1

	x, y := s[0], s[1]

	if x > 0 && g.layout[y][x-1] != "#" {
		adj = append(adj, []int{x - 1, y})
	}

	if x < maxX && g.layout[y][x+1] != "#" {
		adj = append(adj, []int{x + 1, y})
	}

	if y > 0 && g.layout[y-1][x] != "#" {
		adj = append(adj, []int{x, y - 1})
	}

	if y < maxY && g.layout[y+1][x] != "#" {
		adj = append(adj, []int{x, y + 1})
	}

	return adj
}

func countOccurence(char string, s []string) int {
	counter := 0
	for _, c := range s {
		if c == char {
			counter++
		}
	}

	return counter
}

func compress(input [][]int) (map[int]int, map[int]int, int, int) {
	visitedX := map[int]struct{}{}
	visitedY := map[int]struct{}{}
	allX := []int{}
	allY := []int{}
	rows := 0
	columns := 0

	for i := range input {
		x, y := input[i][0], input[i][1]
		// unique
		if _, ok := visitedX[x]; !ok {
			allX = append(allX, x)
		}
		if _, ok := visitedY[y]; !ok {
			allY = append(allY, y)
		}

		visitedX[x] = struct{}{}
		visitedY[y] = struct{}{}
	}

	// sort
	slices.SortFunc(allX, func(a, b int) int {
		return cmp.Compare(a, b)
	})
	// sort
	slices.SortFunc(allY, func(a, b int) int {
		return cmp.Compare(a, b)
	})

	compressedX := map[int]int{}
	compressedY := map[int]int{}
	for i := range len(allX) {
		compressedX[allX[i]] = i
	}

	for i := range len(allY) {
		compressedY[allY[i]] = i
	}

	columns = len(allX)
	rows = len(allY)

	return compressedX, compressedY, rows, columns
}

// return coordinates which have to be filled
func rasterize(input [][]int) [][]int {
	var out [][]int
	for i := range len(input) - 1 {
		a := input[i]
		b := input[i+1%len(input)]

		aX, bX := a[0], b[0]
		aY, bY := a[1], b[1]

		// When both X are the same, then Y is different
		// and vice versa
		if aX == bX {
			y0, y1 := min(aY, bY), max(aY, bY)
			for y := y0 + 1; y < y1; y++ {
				out = append(out, []int{aX, y})
			}
		} else if aY == bY {
			x0, x1 := min(aX, bX), max(aX, bX)
			for x := x0 + 1; x < x1; x++ {
				out = append(out, []int{x, aY})
			}
		}
	}

	return out
}

func findLargestArea(input [][]int) int {

	largestArea := 0
	for i := range input {

		x1, y1 := input[i][0], input[i][1]

		for j := i + 1; j < len(input); j++ {

			x2, y2 := input[j][0], input[j][1]

			s := square{
				width:  max(x1, x2) - min(x1, x2) + 1,
				heigth: max(y1, y2) - min(y1, y2) + 1,
			}
			a := s.area()
			if a > largestArea {
				largestArea = a
			}
		}
	}

	return largestArea
}

type square struct {
	width  int
	heigth int

	a []int
	b []int
}

func (s square) area() int {
	return s.width * s.heigth
}

func partTwo(input [][]int) {
	grid := newGrid(input)
	fmt.Println(grid)

	largestArea := 0
	for i := range len(grid.vertex) - 1 {
		x1, y1 := grid.vertex[i][0], grid.vertex[i][1]

		for j := i + 1; j < len(grid.vertex); j++ {
			x2, y2 := grid.vertex[j][0], grid.vertex[j][1]

			squareVertex := [][]int{
				{max(x1, x2), max(y1, y2)},
				{min(x1, x2), max(y1, y2)},
				{max(x1, x2), min(y1, y2)},
				{min(x1, x2), min(y1, y2)},
			}
			o1 := grid.lookupOriginalCoordinate(x1, y1)
			o2 := grid.lookupOriginalCoordinate(x2, y2)

			oX1, oY1 := o1[0], o1[1]
			oX2, oY2 := o2[0], o2[1]

			s := square{
				width:  max(oX1, oX2) - min(oX1, oX2) + 1,
				heigth: max(oY1, oY2) - min(oY1, oY2) + 1,
			}
			a := s.area()

			if grid.insidePolygon(squareVertex) {
				// fmt.Println(x1, y1, "->", x2, y2, s.area(), "(", oX1, oY1, "->", oX2, oY2, ")")
				if a > largestArea {
					largestArea = a
				}
			}
		}
	}

	fmt.Println(largestArea)
}

func main() {
	file, err := os.Open("inputSample.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	puzzleInput := [][]int{}
	for scanner.Scan() {

		line := strings.Split(scanner.Text(), ",")
		col, row := line[0], line[1]

		colInt, err := strconv.Atoi(col)
		if err != nil {
			log.Fatal("col", err)
		}

		rowInt, err := strconv.Atoi(row)
		if err != nil {
			log.Fatal("row", err)
		}

		puzzleInput = append(puzzleInput, []int{colInt, rowInt})
	}

	largest := findLargestArea(puzzleInput)
	fmt.Println(largest)

	partTwo(puzzleInput)
}
