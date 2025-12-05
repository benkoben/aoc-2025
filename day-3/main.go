package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type bank struct {
	batteries []int
}

func newBank(batteries string) (*bank, error) {
	out := &bank{}

	for _, b := range batteries {
		batteryNum, err := strconv.Atoi(string(b))
		if err != nil {
			return nil, err
		}
		out.batteries = append(out.batteries, batteryNum)
	}

	return out, nil
}

func (b bank) maxJolt() int {
	var maxJolt int
	var first int
	var second int
	var atIndex int

	// Find the max number in the bank
	first, atIndex = maxIndex(b.batteries)
	// If the largest number is at the end, we have found the second number
	// Switch and contiue to find the first number
	if atIndex == len(b.batteries)-1 {
		second = first
		first, atIndex = maxIndex(b.batteries[:atIndex])
	} else {
		// If we have found the first number (the largest not at the end)
		// continue then we need to find the second largest number (which comes after the first in the slice)
		second, atIndex = maxIndex(b.batteries[atIndex+1:])
	}
	// fmt.Println(first, second, atIndex, len(b.batteries))

	maxJolt, _ = strconv.Atoi(fmt.Sprintf("%d%d", first, second))

	return maxJolt
}

func maxIndex(in []int) (int, int) {
	var max int
	var index int

	for i, num := range in {
		if num > max {
			max = num
			index = i
		}
	}

	return max, index
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	banks := []bank{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()

		bank, err := newBank(input)
		if err != nil {
			log.Fatalln(err)
		}
		banks = append(banks, *bank)
	}

	resultPartOne := 0
	for _, bank := range banks {
		// fmt.Println(bank.batteries, bank.maxJolt())
		resultPartOne += bank.maxJolt()
	}

	fmt.Println("result part one: ", resultPartOne)
}
