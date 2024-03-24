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

	var memStatsList []MemStatsAvg
	for _, record := range records {
		var memStats MemStatsAvg
		memStats.N, _ = strconv.Atoi(record[0])
		memStats.Time, _ = strconv.ParseFloat(record[1], 64)
		memStats.Alloc, _ = strconv.ParseFloat(record[2], 64)
		memStats.TotalAlloc, _ = strconv.ParseFloat(record[3], 64)
		memStats.Sys, _ = strconv.ParseFloat(record[4], 64)
		memStats.Lookups, _ = strconv.ParseFloat(record[5], 64)
		memStats.Mallocs, _ = strconv.ParseFloat(record[6], 64)
		memStats.Frees, _ = strconv.ParseFloat(record[7], 64)
		memStats.HeapAlloc, _ = strconv.ParseFloat(record[8], 64)
		memStats.HeapSys, _ = strconv.ParseFloat(record[9], 64)
		memStats.HeapIdle, _ = strconv.ParseFloat(record[10], 64)
		memStats.HeapInuse, _ = strconv.ParseFloat(record[11], 64)
		memStats.HeapReleased, _ = strconv.ParseFloat(record[12], 64)
		memStats.HeapObjects, _ = strconv.ParseFloat(record[13], 64)
		memStats.StackInuse, _ = strconv.ParseFloat(record[14], 64)
		memStats.StackSys, _ = strconv.ParseFloat(record[15], 64)
		memStats.MSpanInuse, _ = strconv.ParseFloat(record[16], 64)
		memStats.MSpanSys, _ = strconv.ParseFloat(record[17], 64)
		memStats.MCacheInuse, _ = strconv.ParseFloat(record[18], 64)
		memStats.MCacheSys, _ = strconv.ParseFloat(record[19], 64)
		memStats.BuckHashSys, _ = strconv.ParseFloat(record[20], 64)
		memStats.GCSys, _ = strconv.ParseFloat(record[21], 64)
		memStats.OtherSys, _ = strconv.ParseFloat(record[22], 64)
		memStats.NextGC, _ = strconv.ParseFloat(record[23], 64)
		memStats.LastGC, _ = strconv.ParseFloat(record[24], 64)

		memStatsList = append(memStatsList, memStats)
	}

	var dataPoints []opts.Chart3DData

	for i, memStats := range memStatsList {
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{0, i, memStats.N}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{1, i, memStats.Time}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{2, i, memStats.Alloc}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{3, i, memStats.TotalAlloc}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{4, i, memStats.Sys}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{5, i, memStats.Lookups}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{6, i, memStats.Mallocs}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{7, i, memStats.Frees}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{8, i, memStats.HeapAlloc}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{9, i, memStats.HeapSys}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{10, i, memStats.HeapIdle}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{11, i, memStats.HeapInuse}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{12, i, memStats.HeapReleased}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{13, i, memStats.HeapObjects}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{14, i, memStats.StackInuse}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{15, i, memStats.StackSys}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{16, i, memStats.MSpanInuse}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{17, i, memStats.MSpanSys}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{18, i, memStats.MCacheInuse}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{19, i, memStats.MCacheSys}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{20, i, memStats.BuckHashSys}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{21, i, memStats.GCSys}})
		dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{22, i, memStats.OtherSys}})
		// Exlude these values, as they often distort the plots to the point of making them useless
		//dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{23, i, memStats.NextGC}})
		//dataPoints = append(dataPoints, opts.Chart3DData{Value: []interface{}{24, i, memStats.LastGC}})
	}

	bar3d := charts.NewBar3D()
	bar3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "MemStats 3D Bar Chart",
			Subtitle: "Visualization of Memory Statistics",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Max:        100, // Adjust based on your data
			InRange:    &opts.VisualMapInRange{Color: []string{"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8", "#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026"}},
		}),
		charts.WithGrid3DOpts(opts.Grid3D{
			BoxWidth:  200,
			BoxDepth:  80,
			BoxHeight: 50,
		}),
	)

	// Setting X and Y axis data
	xAxisData := []string{
		"N", "Time", "Alloc", "TotalAlloc", "Sys", "Lookups", "Mallocs", "Frees", "HeapAlloc", "HeapSys",
		"HeapIdle", "HeapInuse", "HeapReleased", "HeapObjects", "StackInuse", "StackSys", "MSpanInuse",
		"MSpanSys", "MCacheInuse", "MCacheSys", "BuckHashSys", "GCSys", "OtherSys",
		/*"NextGC", "LastGC"*/ // Exlude these values, as they often distort the plots to the point of making them useless
	}
	yAxisData := make([]string, len(memStatsList))
	for i := range yAxisData {
		yAxisData[i] = strconv.Itoa(i)
	}

	// Adding series with X and Y axis data
	bar3d.AddSeries("bar3d", dataPoints).
		SetGlobalOptions(
			charts.WithXAxis3DOpts(opts.XAxis3D{Data: xAxisData}),
			charts.WithYAxis3DOpts(opts.YAxis3D{Data: yAxisData}),
		)

	// Save the chart to a file
	f, err := os.Create("bar3d.html")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	if err := bar3d.Render(f); err != nil {
		log.Fatalf("Failed to render chart: %s", err)
	}

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
			memStats[fieldName] = append(memStats[fieldName], opts.LineData{Value: v.Field(i).Interface()})
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
	line.Render(f)
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
	plotStatsBar3D()
	//plotStatsLine()

}
