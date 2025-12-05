package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	start int = 50
)

type step struct {
	direction string
	distance  int
}

type safe struct {
	position int
	password int
}

func newSafe() safe {
	return safe{
		position: start,
		password: 0,
	}
}

func (s *safe) turnLeft(distance int) {
	// Save the revolutions and remove them from the distance
	revolutions := int(math.Abs(float64(distance) / 100))
	if revolutions > 0 {
		distance -= 100 * revolutions

		// Adjust password with revolutions
		s.password += revolutions
	}
	// Calculate the new position
	newPos := ((s.position - distance) + 100) % 100

	// Check the conditions when the newPosition lands on 0
	// Or when we 0 is passed
	if newPos == 0 {
		s.password += 1
	} else if s.position > 0 && s.position-distance < 0 {
		s.password += 1
	}

	oldPos := s.position
	// Update the position
	s.position = newPos
	fmt.Printf("%d -> L(%d)%d -> %d\n", oldPos, revolutions, distance, s.position)
}

func (s *safe) turnRight(distance int) {
	// Save the revolutions and remove them from the distance
	revolutions := int(math.Abs(float64(distance) / 100))
	if revolutions > 0 {
		distance -= 100 * revolutions

		// Adjust password with revolutions
		s.password += revolutions
	}
	newPos := ((s.position + distance) + 100) % 100

	// Check the conditions when the newPosition lands on 0
	// Or when we 0 is passed
	if newPos == 0 {
		s.password += 1
	} else if s.position > 0 && s.position+distance > 100 {
		s.password += 1
	}

	oldPos := s.position

	// Update the position
	s.position = newPos
	fmt.Printf("%d -> R(%d)%d -> %d\n", oldPos, revolutions, distance, s.position)
}

func parseInput(in []string) []step {
	out := make([]step, len(in))
	for i := range in {

		direction := string(in[i][0])
		distance := string(in[i][1:])

		distanceInt, err := strconv.Atoi(distance)
		if err != nil {
			log.Fatalln(err)
		}

		out[i] = step{direction, distanceInt}
	}

	return out
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	puzzleInput := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		puzzleInput = append(puzzleInput, scanner.Text())
	}

	in := parseInput(puzzleInput)
	safe := newSafe()
	for _, step := range in {
		switch step.direction {
		case "R":
			safe.turnRight(step.distance)
		case "L":
			safe.turnLeft(step.distance)
		}
	}

	fmt.Printf("password: %d\n", safe.password)
}
