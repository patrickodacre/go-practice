package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"flag"
)

func main() {

	var (
		problemSetFile = flag.String("problemset", "problems.csv", "Specify a problem set file.")
	)

	flag.Parse()

	problemSet, err := os.OpenFile(*problemSetFile, os.O_RDWR, 0755)

	failOnError(err, "Error opening problem set.")

	defer problemSet.Close()

	r := csv.NewReader(problemSet)

	// counters
	numOfCorrectAnswers := 0
	numOfQuestions := 0

	resultsReport := [][]string{}
	now, _ := time.Now().MarshalText()
	timestamp := string(string(now)[:10])

	resultsReport = append(resultsReport, []string{"Results:", timestamp})

	for {
		reportEntry := []string{}

		record, err := r.Read()

		if err == io.EOF {
			break
		}

		failOnError(err, "Error reading problem set.")

		question := record[0]
		expectedAnswer := record[1]

		numOfQuestions++

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

		// fmt.Println(string(response))
		fmt.Println("...saved")

		answerGiven = string(answerGiven)
		answerGiven = strings.Trim(answerGiven, " ")
		answerGiven = answerGiven[:len(answerGiven)-1] // trim off the \n

		// record the answer given w/o the \n or whitespace
		reportEntry = append(reportEntry, answerGiven)

		// verify the answer is a number:
		_, err = strconv.Atoi(answerGiven)

		// invalid input
		if err != nil {
			reportEntry = append(reportEntry, "false")

			continue
		}

		if answerGiven == expectedAnswer {
			numOfCorrectAnswers++
			reportEntry = append(reportEntry, "true")
		} else {
			reportEntry = append(reportEntry, "false")
		}

		resultsReport = append(resultsReport, reportEntry)
	}

	resultsReport = append(resultsReport, []string{"", "", "Final Score", strconv.Itoa(numOfCorrectAnswers) + "/" + strconv.Itoa(numOfQuestions)})

	reportFile, err := os.Create(timestamp + "-report.csv")

	failOnError(err, "Failed to generate report.")

	defer reportFile.Close()

	reportFile.Chmod(0755)

	w := csv.NewWriter(reportFile)

	for _, record := range resultsReport {
		if err := w.Write(record); err != nil {
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
