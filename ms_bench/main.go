package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Open the Excel file
	f, err := excelize.OpenFile("DataBenchmark.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	// Read usernames and passwords from the Excel file
	rows, err := f.GetRows("Sheet1") // Assuming the data is in Sheet1
	if err != nil {
		log.Fatal(err)
	}

	var targets []vegeta.Target
	for index, row := range rows {
		if index == 0 {
			continue
		} else if index == 200 {
			break
		}
		// Assuming the first column contains usernames and the second column contains passwords
		if len(row) >= 2 {
			target := vegeta.Target{
				Method: "POST",
				URL:    "http://localhost:8080/user/login", // Adjust URL as needed
				Body:   []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, row[2], row[3])),
			}
			targets = append(targets, target)
		}
	}

	// Create a new Vegeta targeter with the targets
	targeter := vegeta.NewStaticTargeter(targets...)

	// Create a new Vegeta attacker
	attacker := vegeta.NewAttacker()

	// Create a new rate limiter to control the request rate (adjust as needed)
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}

	// Initialize metrics
	var metrics vegeta.Metrics

	// Create a channel to receive SIGINT and SIGTERM signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start attacking in a separate goroutine
	go func() {
		for res := range attacker.Attack(targeter, rate, 10*time.Second, "Load test") {
			metrics.Add(res)
		}
	}()

	// Wait for SIGINT or SIGTERM to stop the attack
	<-interrupt

	// Stop the attacker
	attacker.Stop()

	// Output attack metrics
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Rate: %f\n", metrics.Rate)
	fmt.Printf("Success: %f%%\n", metrics.Success*100)
	fmt.Printf("Status Codes: %v\n", metrics.StatusCodes)
	fmt.Printf("Errors: %d\n", metrics.Errors)

	// Output attack latency distribution
	fmt.Println("Latencies:")
	latencies := metrics.Latencies
	fmt.Printf("Mean: %s\n", latencies.Mean)
	fmt.Printf("50th percentile: %s\n", latencies.P50)
	fmt.Printf("95th percentile: %s\n", latencies.P95)
	fmt.Printf("99th percentile: %s\n", latencies.P99)
}
