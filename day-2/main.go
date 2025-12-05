package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type IDRange struct {
	start int
	stop  int
}

func (id IDRange) validate() []int {
	invalidNumbers := []int{}

	for i := id.start; i <= id.stop; i++ {

		iStr := strconv.Itoa(i)
		if len(iStr)%2 == 0 {

			chunk1Start, chunk1Stop := 0, (len(iStr) / 2)
			chunk2Start, chunk2Stop := len(iStr)/2, len(iStr)

			if iStr[chunk1Start:chunk1Stop] == iStr[chunk2Start:chunk2Stop] {
				invalidNumbers = append(invalidNumbers, i)
			}
		}
	}

	return invalidNumbers
}

func (id IDRange) validate2() []int {
	invalidNumbers := []int{}

	// Find invalid numbers by comparing two chunks.
	// ChunkOne always stays at the beginning.
	// ChunkTwo is put after chunkOne
	// If there are no matches between chunkOne and chunkTwo then increase the chunk sizes.
	// If there is a match then:
	// 1. If chunkTwo has reached the end, then break the loop <- we have found an invalid number
	// 2. If ChunkTwo has not reached the end, move it one slot closer to end end and re-evaluate.
	// 3. If no match is found, both cells grow in size.
	for i := id.start; i <= id.stop; i++ {
		iStr := strconv.Itoa(i)
		chunkStart := 0
		chunkStop := 0
		chunkSize := 1
		chunkTwoStart := chunkStart + chunkSize
		chunkTwoStop := chunkTwoStart + chunkSize
		for chunkStop <= len(iStr) {

			if chunkTwoStop > len(iStr) {
				// Short circuit if chunkTwo's stop index is greater than the input number's length.
				// This means all indexes have been exhausted and we cannot continue.
				break
			}

			chunkStop = chunkStart + chunkSize
			chunkOne := iStr[chunkStart:chunkStop]
			chunkTwo := iStr[chunkTwoStart:chunkTwoStop]

			if chunkOne == chunkTwo && chunkTwoStop == len(iStr) {
				// If the chunks match and the last chunk's is at the end
				// This means the number is made up of **only sequences**.
				invalidNumbers = append(invalidNumbers, i)
				break
			} else if chunkOne == chunkTwo && chunkTwoStop < len(iStr) {
				// the two chunks match but we have not reached the end yet.
				// move chunk 2 into the next possible chunk slot and re-evaluate
				chunkTwoStart = chunkTwoStop
				chunkTwoStop = chunkTwoStop + chunkSize
				continue
			} else {
				// If the chunks dont match keep increasing the chunks sizes.
				chunkTwoStart = chunkStop + 1
				chunkTwoStop = chunkTwoStart + chunkSize + 1
				chunkSize++
			}
		}
	}

	return invalidNumbers
}

func sum(in []int) int {
	total := 0

	for _, i := range in {
		total += i
	}
	return total
}

func main() {

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	resultPartOne := []int{}
	resultPartTwo := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		for productIdRange := range strings.SplitSeq(input, ",") {
			id, err := newIDRange(productIdRange)
			if err != nil {
				log.Fatalln(err)
			}
			resultPartOne = append(resultPartOne, id.validate()...)
			resultPartTwo = append(resultPartTwo, id.validate2()...)

		}
	}

	fmt.Println("result part one: ", sum(resultPartOne))
	fmt.Println("result part two: ", sum(resultPartTwo))

}

func newIDRange(in string) (*IDRange, error) {
	productIdRange := strings.Split(in, "-")

	start, err := strconv.Atoi(productIdRange[0])
	if err != nil {
		return nil, err
	}

	stop, err := strconv.Atoi(productIdRange[1])
	if err != nil {
		return nil, err
	}

	return &IDRange{
		start: start,
		stop:  stop,
	}, nil
}
