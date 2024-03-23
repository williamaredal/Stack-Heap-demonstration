package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

func recursiveFunction(count int) int {
	if count <= 0 {
		return 0
	}
	return 1 + recursiveFunction(count-1)
}

func loopFunction(count int) int {
	sum := 0
	for i := 0; i < count; i++ {
		sum += 1
	}

	return sum
}

func writeMemStatsToCSV(writeCount int, n int, elapsedTime int64) {
	var file *os.File
	var err error
	headers := []string{
		"N", "Time", "Alloc", "TotalAlloc", "Sys", "Lookups", "Mallocs", "Frees",
		"HeapAlloc", "HeapSys", "HeapIdle", "HeapInuse", "HeapReleased",
		"HeapObjects", "StackInuse", "StackSys", "MSpanInuse", "MSpanSys",
		"MCacheInuse", "MCacheSys", "BuckHashSys", "GCSys", "OtherSys",
		"NextGC", "LastGC",
	}

	if writeCount == 0 {
		// Create a new file or truncate an existing one
		file, err = os.Create("memstats.csv")
	} else {
		// Append to an existing file
		file, err = os.OpenFile("memstats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	if err != nil {
		log.Fatalf("failed creating or opening file: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if writeCount == 0 {
		// Write headers only if it's a new file
		if err := writer.Write(headers); err != nil {
			log.Fatalln("error writing headers to csv:", err)
		}
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := []string{
		fmt.Sprintf("%d", n),
		fmt.Sprintf("%d", elapsedTime),
		fmt.Sprintf("%d", m.Alloc),
		fmt.Sprintf("%d", m.TotalAlloc),
		fmt.Sprintf("%d", m.Sys),
		fmt.Sprintf("%d", m.Lookups),
		fmt.Sprintf("%d", m.Mallocs),
		fmt.Sprintf("%d", m.Frees),
		fmt.Sprintf("%d", m.HeapAlloc),
		fmt.Sprintf("%d", m.HeapSys),
		fmt.Sprintf("%d", m.HeapIdle),
		fmt.Sprintf("%d", m.HeapInuse),
		fmt.Sprintf("%d", m.HeapReleased),
		fmt.Sprintf("%d", m.HeapObjects),
		fmt.Sprintf("%d", m.StackInuse),
		fmt.Sprintf("%d", m.StackSys),
		fmt.Sprintf("%d", m.MSpanInuse),
		fmt.Sprintf("%d", m.MSpanSys),
		fmt.Sprintf("%d", m.MCacheInuse),
		fmt.Sprintf("%d", m.MCacheSys),
		fmt.Sprintf("%d", m.BuckHashSys),
		fmt.Sprintf("%d", m.GCSys),
		fmt.Sprintf("%d", m.OtherSys),
		fmt.Sprintf("%d", m.NextGC),
		fmt.Sprintf("%d", m.LastGC),
	}

	if err := writer.Write(stats); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
}

func main() {
	// Turns off automatic garbage collection before tests
	debug.SetGCPercent(-1)

	// Write count to know if headers should be written, or rows appended
	writeCounter := 0

	// Element depth to test
	for d := 0; d < 7; d++ {
		// Run the function with n elements
		n := int(10 * math.Pow(10, float64(d)))

		// Number of times to test any given element depth
		for test_i := 0; test_i < 10; test_i++ {
			startTime := time.Now()
			//result := recursiveFunction(int(n))
			result := loopFunction(n)
			elapsedTime := time.Since(startTime)

			// Call the function to write mem stats to CSV
			writeMemStatsToCSV(writeCounter, n, elapsedTime.Milliseconds())

			// Displays the results
			fmt.Printf("Current runtime: %dms\nRecursive depth: %d\nFunction result: %d\n\n", elapsedTime, int(n), result)

			// Runs garbage collection to clear mem-usage before next run
			runtime.GC()

			// Increments writes after test is finished
			writeCounter++
		}

	}
}
