// Package main ....
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// execute_quiz_question prompts the user with a question and awaits user
// input; once received, a boolean is returned indicating whether the response
// matches the solution.
func execute_quiz_question(
	question_num int,
	question string,
	solution string,
) bool {
	fmt.Printf("%v. %v?\n", question_num, question)

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	return response == solution
}

// main reads in data, provided via a CSV file, and executes a timed quiz.
func main() {
	shuffleFlag := flag.Bool(
		"shuffle",
		false,
		"Indicates whether or not to shuffle the quiz; default is `false`.",
	)
	time_limitFlag := flag.Int(
		"time_limit",
		30,
		"Sets the time duration of the quiz (in seconds); default is 30",
	)

	f, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	shuffle := *shuffleFlag
	time_limit := *time_limitFlag

	if shuffle {
		fmt.Print("Shuffling is enabled.\nShuffling quiz....\n")
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(
			len(data),
			func(i, j int) {
				data[i], data[j] = data[j], data[i]
			},
		)
		fmt.Println("Shuffling successful.")
	}

	fmt.Printf(
		"You will have %v seconds to answer %v questions.\n",
		time_limit,
		len(data),
	)
	fmt.Println("Please press enter to begin the quiz.")
	fmt.Scanln()

	timer := time.NewTimer(time.Duration(time_limit) * time.Second)

	total := 0
	for i, row := range data {

		result := make(chan bool, 1)
		go func() {
			result <- execute_quiz_question(i, row[0], row[1])
		}()

		select {
		case <-timer.C:
			fmt.Println("Time's up!")
			fmt.Printf(
				"You got %v questions right out of %v.\n", total, len(data),
			)
			return
		case val := <-result:
			if val {
				total++
			}
		}
	}

	fmt.Printf("You got %v questions right out of %v.\n", total, len(data))
}
