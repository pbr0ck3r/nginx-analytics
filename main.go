package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
)

type Log struct {
	Time       string `json:"time"`
	RemoteIp   string `json:"remote_ip"`
	RemoteUser string `json:"remote_user"`
	Request    string `json:"request"`
	Response   int    `json:"response"`
	Bytes      int    `json:"bytes"`
	Referrer   string `json:"referrer"`
	Agent      string `json:"agent"`
}

func main() {
	statusCodes := make(map[int]int)
	errorCounts := make(map[string]int)
	largestRequestSize := 0
	var largestRequest Log
	allSizes := make([]int64, 0)
	successSizes := make([]int64, 0)
	errorSizes := make([]int64, 0)

	// Open the file
	file, err := os.Open("nginx_json_logs.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read each line of the file
	for scanner.Scan() {
		// Unmarshal each line as JSON
		var logEntry Log
		err := json.Unmarshal(scanner.Bytes(), &logEntry)
		if err != nil {
			fmt.Printf("Error parsing line: %s\n", scanner.Text())
			continue // Skip to next line on error
		}

		// Calculate count of each status code
		statusCodes[logEntry.Response]++

		// Prepare for calculating mean/median/p99 response bytes
		allSizes = append(allSizes, int64(logEntry.Bytes))
		if logEntry.Response < 400 {
			// Success
			successSizes = append(successSizes, int64(logEntry.Bytes))
		} else {
			// error
			errorSizes = append(errorSizes, int64(logEntry.Bytes))
		}

		// Determine which endpoint returned the largest Response by bytes
		if logEntry.Bytes > largestRequestSize {
			largestRequestSize = int(logEntry.Bytes)
			largestRequest = logEntry
		}

		// Prepare for determining which endpoint returned the most error repsonses (>= 400 code)
		if logEntry.Response >= 400 {
			errorCounts[logEntry.RemoteIp]++
		}
	}

	fmt.Println("What is the count of each status code?")
	for statusCode, count := range statusCodes {
		fmt.Printf("Status code %d: %d\n", statusCode, count)
	}
	fmt.Println()

	fmt.Println("What is the mean/median/p99 response bytes?")
	fmt.Printf("All Requests: Mean: %.2f, Median: %d, p99: %d\n", calculateMean(allSizes), calculateMedian(allSizes), calculateP99(allSizes))
	fmt.Printf("All Successful Requests (< 400): Mean: %.2f, Median: %d, p99: %d\n", calculateMean(successSizes), calculateMedian(successSizes), calculateP99(successSizes))
	fmt.Printf("All Error Requests (>= 400): Mean: %.2f, Median: %d, p99: %d\n", calculateMean(errorSizes), calculateMedian(errorSizes), calculateP99(errorSizes))
	fmt.Println()

	fmt.Println("Which endpoint returned the largest Response by bytes?")
	fmt.Printf("%s (%d bytes)\n", largestRequest.RemoteIp, largestRequestSize)
	fmt.Println()

	fmt.Println("Which endpoint returned the most error repsonses (>= 400 code)?")
	mostErrors := 0
	var mostErrorEndpoint string
	for endpoint, count := range errorCounts {
		if count > mostErrors {
			mostErrors = count
			mostErrorEndpoint = endpoint
		}
	}
	fmt.Printf("%s (%d errors)\n", mostErrorEndpoint, mostErrors)
}

func calculateMean(data []int64) float64 {
	var sum float64
	for _, d := range data {
		sum += float64(d)
	}
	return sum / float64(len(data))
}

func calculateMedian(data []int64) int64 {
	// Ensure data is sorted
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })
	mid := len(data) / 2
	if len(data)%2 == 0 {
		// Even number of elements, median is average of middle two
		return (data[mid-1] + data[mid]) / 2
	}
	// Odd number of elements, median is the middle element
	return data[mid]
}

func calculateP99(data []int64) int64 {
	if len(data) == 0 {
		return 0
	}
	// Ensure data is sorted
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})
	index := int(math.Round(0.99 * float64(len(data))))
	return data[index]
}
