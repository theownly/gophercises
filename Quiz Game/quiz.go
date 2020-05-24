package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	var csvFilePath = flag.String("csv", "problems.csv", "relative path to csv file")
	var timeLimit = flag.Int("time", 30, "time limit in seconds")
	flag.Parse()

	csvFile, err := os.Open(*csvFilePath)

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()

	csvFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	type problem struct {
		question string
		solution string
	}

	problems := make([]problem, len(records))

	for i, record := range records {
		problems[i] = problem{
			question: record[0],
			solution: record[1],
		}
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	inputChannel := make(chan string)

	type counter struct {
		total   int
		correct int
		wrong   int
	}

	c := counter{
		total: len(problems),
	}

	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		go func() {
			var input string
			fmt.Scanf("%s\n", &input)
			inputChannel <- input
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTimes up! You answered %d out of %d problems. Out of these you had %d right and %d wrong.\n",
				c.correct+c.wrong, c.total, c.correct, c.wrong)
			return
		case input := <-inputChannel:
			if input == problem.solution {
				c.correct++
			} else {
				c.wrong++
			}
		}
	}

	fmt.Printf("\nGreat! You answered all %d problems. Out of these you had %d right and %d wrong.\n",
		c.total, c.correct, c.wrong)

}
