package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type cell string

func (c cell) isPaperRoll() bool {
	return c == "@"
}

type grid struct {
	rows   []row
	top    int
	bottom int
}

type row struct {
	columns []cell
	left    int
	right   int

	// Added in Part 2
	paperRollsCount int
}

func newGridFromInput(input []string) grid {
	rows := make([]row, len(input))
	for i := range input {
		rows[i] = convertToRow(input[i])
	}

	return grid{
		rows:   rows,
		top:    0,
		bottom: len(rows) - 1,
	}
}

// findMovablePaperRolls walk the grid and mark all cells that have less than threshold amount of paper rolls within adjacent cells
func (g *grid) findMovablePaperRolls(threshold int) int {
	movableCount := 0

	for r := g.top; r <= g.bottom; r++ {
		row := g.rows[r]

		// Added in part 2
		if row.paperRollsCount == 0 {
			// Short circuit if the row has no paper rolls
			continue
		}

		for col := row.left; col <= row.right; col++ {
			cell := row.columns[col]
			if !cell.isPaperRoll() {
				// Short circuit if the cell is not a paper roll
				continue
			}

			paperRollCount := 0
			canPeekLeft := col > row.left
			canPeekRight := col < row.right
			canPeekUp := r > g.top
			canPeekDown := r < g.bottom
			canPeekUpLeft := canPeekLeft && canPeekUp
			canPeekDownLeft := canPeekLeft && canPeekDown
			canPeekUpRight := canPeekRight && canPeekUp
			canPeekDownRight := canPeekRight && canPeekDown

			if canPeekLeft {
				leftNeighbor := row.columns[col-1]
				if leftNeighbor.isPaperRoll() {
					paperRollCount++
				}
			}

			if canPeekRight {
				rightNeighbor := row.columns[col+1]
				if rightNeighbor.isPaperRoll() {
					paperRollCount++
				}
			}

			if canPeekUp {
				upNeighbor := g.rows[r-1].columns[col]
				if upNeighbor.isPaperRoll() {
					paperRollCount++
				}
			}

			if canPeekDown {
				downNeighbor := g.rows[r+1].columns[col]
				if downNeighbor.isPaperRoll() {
					paperRollCount++
				}
			}

			if canPeekUpLeft {
				leftUpNeighbor := g.rows[r-1].columns[col-1]
				if leftUpNeighbor.isPaperRoll() {
					paperRollCount++
				}
			}

			if canPeekDownLeft {
				leftDownNeighbor := g.rows[r+1].columns[col-1]
				if leftDownNeighbor.isPaperRoll() {
					paperRollCount++
				}
			}

			if canPeekUpRight {
				rightUpNeighbor := g.rows[r-1].columns[col+1]
				if rightUpNeighbor.isPaperRoll() {
					paperRollCount++
				}
			}

			if canPeekDownRight {
				rightDownNeighbor := g.rows[r+1].columns[col+1]
				if rightDownNeighbor.isPaperRoll() {
					paperRollCount++
				}
			}

			if paperRollCount < threshold {
				// Record
				movableCount++

				// Added in Part 2
				// Convert
				g.rows[r].columns[col] = "."
				// Count down number of paper rolls
				g.rows[r].paperRollsCount--
			}
		}
	}

	return movableCount
}

// convertToRow converts a string to a row of cells
// it returns the row and the number of paperRolls in the row
func convertToRow(in string) row {
	paperRollCount := 0
	cells := make([]cell, len(in))
	for i := range in {
		cells[i] = cell(in[i])

		// PART 2
		if cells[i].isPaperRoll() {
			paperRollCount++
		}
	}
	return row{cells, 0, len(cells) - 1, paperRollCount}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	puzzleInput := []string{}
	for scanner.Scan() {
		row := scanner.Text()
		puzzleInput = append(puzzleInput, row)
	}

	newGrid := newGridFromInput(puzzleInput)
	result := 0
	iterationCount := 0
	for {
		movableCount := newGrid.findMovablePaperRolls(4)
		if movableCount == 0 {
			break
		}

		log.Printf("iteration %d -> movable paper rolls: %d\n", iterationCount, movableCount)

		result += movableCount
		iterationCount++
	}
	fmt.Println("result part one: ", result)
}
