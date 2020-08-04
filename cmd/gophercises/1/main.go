package main

import (
	"bufio"
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

	// Read the input to start
	buf := bufio.NewReader(os.Stdin)

	fmt.Println("You will have " + strconv.Itoa(*limitSeconds) + " seconds to complete the quiz.")
	fmt.Println("Hit enter when you're ready to start...")

	// ReadBytes until \n which is "enter".
	// because \n is our delimeter, we have to trim
	// that off our answer to do a comparison with the
	// expected answer.
	_, err = buf.ReadBytes('\n')

	failOnError(err, "Failed to start the quiz")

	// countdown
	go func() {

		limit := *limitSeconds

		for range time.Tick(1 * time.Second) {
			limit--

			if limit <= 0 {
				timesup = true
				log.Println("Time's up!")
				break
			}
		}
	}()

	for _, problem := range problems {
		reportEntry := []string{}

		question := problem[0]
		expectedAnswer := strings.ToLower(problem[1])

		// if the time is up, we want to continue to
		// write questions and results to the final report
		if timesup {
			reportEntry = append(reportEntry, question)
			reportEntry = append(reportEntry, "time limit reached")
			reportEntry = append(reportEntry, "--")

			resultsReport = append(resultsReport, reportEntry)
			continue
		}

		// Read the input / answer
		buf := bufio.NewReader(os.Stdin)

		reportEntry = append(reportEntry, question)

		fmt.Println("> " + question + " = ?")

		// ReadBytes until \n which is "enter".
		// because \n is our delimeter, we have to trim
		// that off our answer to do a comparison with the
		// expected answer.
		input, err := buf.ReadBytes('\n')

		if err != nil {
			reportEntry = append(reportEntry, "invalid")
			reportEntry = append(reportEntry, "false")

			continue
		}

		answerGiven := string(input)

		fmt.Println("...saved")

		answerGiven = strings.Trim(answerGiven, " ")
		answerGiven = strings.ToLower(answerGiven)
		answerGiven = answerGiven[:len(answerGiven)-1] // trim off the \n

		// record the answer given w/o the \n or whitespace
		reportEntry = append(reportEntry, answerGiven)

		if answerGiven == expectedAnswer {
			numOfCorrectAnswers++
			reportEntry = append(reportEntry, "true")
		} else {
			reportEntry = append(reportEntry, "false")
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
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
