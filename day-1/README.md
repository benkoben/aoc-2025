# Day 1

First day of advent of code. I have been putting quite a lot of effort the last year to improve my programming knowledge. Which makes we quite stoked for this year's advent of code. Without further a do, lets jump right into it.

# The main challenge

I wont cover the assignment into detail as those are already covered on the official page. But the main key take away for today's problem is:

1. We need to keep track of our position of a safe dial, the position is tracked as a number between 0-100.
2. Each row in our input is an intruction that shifts the position either to the left or the right.
3. There are clear boundaries, moving to the left of 0 or to the right of 100 causes the position to start at the opposite end. 
4. Anytime the dial lands on 0 (part 1) and passes 0 (part 2) after a single instruction we increment a number that represents our flag (themed as a password in the assignment).
5. Passing 0 means that we need to account for revolutions of the dial. A single revolution means 100 positions have been passed within a single intruction.


# The code

```go
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
```

To read and parse the input I used a pretty straight forward method how to read data from text files in go. I wont explain that. Instead I want to focus on the assignment itself.

I defined two structs to manage state and convert text introductions into runtime objects. The `step` struct represents a single row in our puzzle input. I used structs instead of a map to serialize the puzzle input because they make the code more readable. We can use a simple switch block in `main` to read the steps and make descisions on them. Then the `safe` struct is used to keep track of the dial's position and the password.

With the safe struct I can easily add two methods to moving the dial to the left or to the right. The core logic is the same in both methods, just the opposite. Here is how I solved it:

1. For each single revolution we know we passed 0 once. A single revolution exhausts all 100 position on the dial. Therefore I devide the distance with 100, I use absolute numbers and round the number downwards using `int()` and `math.Abs()`. I.e. R224 -> 2 revolutions, `R556` -> 5 revolutions. Whatever the number may be, we'll add it to the safe's password field.
2. Next I simplify the distance by removing the revolution factor, allowing me to calculate the new position solely based lower and upper limits of the dial. A number between 0 and 99. A modulus operator makes this really straight forward.
3. Then we check if our new position lands on 0 or passes 0. In either way we need to increment the safe password.
4. Finally we update the position

