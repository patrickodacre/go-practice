package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var timestamp string
var numOfCorrectAnswers int
var numOfQuestions int
var timesup bool

func main() {

	var (
		problemSetFile = flag.String("problemset", "problems.csv", "Specify a problem set file.")
		limitSeconds   = flag.Int("limitseconds", 30, "Specific a time limit in seconds.")
	)

	flag.Parse()

	problemSet, err := os.OpenFile(*problemSetFile, os.O_RDWR, 0755)

	failOnError(err, "Error opening problem set.")

	defer problemSet.Close()

	r := csv.NewReader(problemSet)

	problems := [][]string{}

	for {
		problem, err := r.Read()

		if err == io.EOF {
			break
		}

		problems = append(problems, problem)
	}

	numOfQuestions = len(problems)

	resultsReport := [][]string{}
	now, _ := time.Now().MarshalText()
	timestamp = string(string(now)[:10])

	resultsReport = append(resultsReport, []string{"Results:", timestamp})

	fmt.Println("You will have " + strconv.Itoa(*limitSeconds) + " seconds to complete the quiz.")
	fmt.Println("Hit enter when you're ready to start...")

	// Hold until ENTER
	fmt.Scanf("\n")

	timer := time.NewTimer(time.Duration(*limitSeconds) * time.Second)

problemsLoop:
	for i, problem := range problems {

		question := problem[0]
		expectedAnswer := strings.ToLower(problem[1])
		reportEntry := []string{question}

		fmt.Printf("Problem %d out of %d : %s = ?", i+1, numOfQuestions, question)
		// create a new channel to accept the answer for the question:
		answerChan := make(chan string)

		go func() {

			var answer string
			// accepts a single string and then ENTER
			fmt.Scanf("%s\n", &answer)

			answer = strings.Trim(answer, " ")
			answer = strings.ToLower(answer)

			answerChan <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("")
			fmt.Println("Time's up!")
			fmt.Println("")
			break problemsLoop
		case answer := <-answerChan:

			fmt.Println("...saved")

			if answer == expectedAnswer {
				numOfCorrectAnswers++

				reportEntry = append(reportEntry, answer, "true")
			} else {
				reportEntry = append(reportEntry, answer, "false")
			}
		}

		resultsReport = append(resultsReport, reportEntry)
	}

	printResults(resultsReport)
}

func printResults(resultsReport [][]string) {

	resultsReport = append(resultsReport, []string{"", "", "Final Score", strconv.Itoa(numOfCorrectAnswers) + "/" + strconv.Itoa(numOfQuestions)})

	reportFile, err := os.Create(timestamp + "-report.csv")

	failOnError(err, "Failed to generate report.")

	defer reportFile.Close()

	reportFile.Chmod(0755)

	w := csv.NewWriter(reportFile)

	for _, result := range resultsReport {
		if err := w.Write(result); err != nil {
			failOnError(err, "Error writing report.")
		}
	}

	w.Flush()

	failOnError(err, "Errors flushing report.")

	fmt.Println("Report Printed")
	os.Exit(1)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
