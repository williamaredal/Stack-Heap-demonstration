package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Struct used for loading and plotting of average stats
type MemStatsAvg struct {
	N            int
	Time         float64
	Alloc        float64
	TotalAlloc   float64
	Sys          float64
	Lookups      float64
	Mallocs      float64
	Frees        float64
	HeapAlloc    float64
	HeapSys      float64
	HeapIdle     float64
	HeapInuse    float64
	HeapReleased float64
	HeapObjects  float64
	StackInuse   float64
	StackSys     float64
	MSpanInuse   float64
	MSpanSys     float64
	MCacheInuse  float64
	MCacheSys    float64
	BuckHashSys  float64
	GCSys        float64
	OtherSys     float64
	NextGC       float64
	LastGC       float64
}

func plotStatsBar3D() {
	file, err := os.Open("memstats.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %s", err)
	}

	records = records[1:] // Skip the header

	memStats := make(map[int][]MemStatsAvg)
	for _, record := range records {
		n, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatalf("Failed to convert N to integer: %s", err)
		}

		var stats MemStatsAvg
		stats.N = n
		if len(record) >= 24 {
			for i, val := range record[1:] {
				floatVal, _ := strconv.ParseFloat(val, 64)
				reflect.ValueOf(&stats).Elem().Field(i + 1).SetFloat(floatVal)
			}
		} else {
			log.Printf("Skipping record due to insufficient fields: %v", record)
		}
		memStats[n] = append(memStats[n], stats)
	}

	bar3D := charts.NewBar3D()
	bar3D.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Memory Stats"}),
		charts.WithXAxis3DOpts(opts.XAxis3D{Name: "N"}),
		charts.WithYAxis3DOpts(opts.YAxis3D{Name: "Stat"}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{Name: "Value"}),
	)

	// Initialize a map to hold all data points
	dataMap := make(map[string][]opts.Chart3DData)

	// Accumulate data points
	for _, statsList := range memStats {
		for _, stats := range statsList {
			v := reflect.ValueOf(stats)
			typeOfS := v.Type()
			for i := 1; i < v.NumField(); i++ {
				fieldName := typeOfS.Field(i).Name
				dataMap[fieldName] = append(dataMap[fieldName], opts.Chart3DData{Value: []interface{}{stats.N, i, v.Field(i).Interface()}})
			}
		}
	}

	// Add series for each stat
	for statName, dataPoints := range dataMap {
		bar3D.AddSeries(statName, dataPoints)
	}

	f, _ := os.Create("bar3d.html")
	bar3D.Render(f)

}

func plotStatsLine() {
	file, err := os.Open("memstats.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %s", err)
	}

	records = records[1:] // Skip the header

	// Using a map where each stat name points to a slice of LineData, holding the N value and the stat value.
	memStats := make(map[string][]opts.LineData)
	for _, record := range records {
		n, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatalf("Failed to convert N to integer: %s", err)
		}

		var stats MemStatsAvg
		stats.N = n
		if len(record) >= 24 {
			for i, val := range record[1:] {
				floatVal, err := strconv.ParseFloat(val, 64)
				if err != nil {
					log.Fatalf("Failed to convert string to float: %s", err)
				}
				reflect.ValueOf(&stats).Elem().Field(i + 1).SetFloat(floatVal)
			}
		} else {
			log.Printf("Skipping record due to insufficient fields: %v", record)
		}

		v := reflect.ValueOf(stats)
		typeOfS := v.Type()
		for i := 1; i < v.NumField(); i++ {
			fieldName := typeOfS.Field(i).Name
			memStats[fieldName] = append(memStats[fieldName], opts.LineData{Value: []interface{}{n, v.Field(i).Interface()}})
		}
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Memory Stats"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "N"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Value"}),
	)

	for statName, dataPoints := range memStats {
		line.AddSeries(statName, dataPoints)
	}

	f, _ := os.Create("line_chart.html")
	if err := line.Render(f); err != nil {
		log.Fatalf("Failed to render chart: %s", err)
	}
}

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
			result := recursiveFunction(int(n))
			//result := loopFunction(n)
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

	debug.SetGCPercent(1)
	// After running the tests, data is loaded and plotted
	//plotStatsBar3D()
	plotStatsLine()

}
