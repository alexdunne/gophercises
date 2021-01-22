package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of question,answer")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	lines := readLines(csvFilename)
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	askQuestions(problems, timer)
}

func readLines(csvFilename *string) [][]string {
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to opent he CSV file: %s\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Couldn't open the file")
	}

	return lines
}

func parseLines(lines [][]string) []Problem {
	problems := make([]Problem, len(lines))

	for i, line := range lines {
		problems[i] = Problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return problems
}

func askQuestions(problems []Problem, timer *time.Timer) {
	correct := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		guessCh := make(chan string)

		go func() {
			// Scanf is a blocking call so instead of having it in the loop
			// we instead move it to a goroutine and send it back via a channel
			// Now in the switch we can wait for either an answer or for the timer to elapse
			var guess string
			fmt.Scanf("%s\n", &guess)

			guessCh <- guess
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTimes up. You scored %d out of %d\n", correct, len(problems))
			return
		case guess := <-guessCh:
			if guess == p.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

type Problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}
