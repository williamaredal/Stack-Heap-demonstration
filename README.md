# Stack-Heap-demonstration
Visualise how a function's resource usage scales with the number of elements (N), across multiple memory metrics to get a better understanding of how your algorithms and datastructures scales. 

This project was initially started after the discovery of how the algorithms and datastructures theory not necessarily translated to the most performant solution ([Bjarne Stroustrup: Why you should avoid Linked Lists](https://www.youtube.com/watch?v=YQs6IC-vgmo&t=1s)).  

Linked List insertion resource usage   |   Array insertion resource usage
:-------------------------:|:-------------------------:
![Illustration Linked List insertion resource growth with N](/Linked-list-insertion.gif)  |  ![Illustration Array insertion resource growth with N](/Array-insertion.gif)


In the example gifs, we can see how the time usage for the linked list grows substantially for each additional N, as compared to the array. In line with what Bjarne Stroustrup noted. The resource usage measurement is not limited to analysis of data structures. Below are two other profiles of functions incrementing a sum by +1, N times. Although the functions are simple, they illustrate the resource footprints of different types of functions. 


Recursive function resource usage   |   For loop function resource usage
:-------------------------:|:-------------------------:
![Illustration Recursive resource growth with N](/Recursive-count.gif)  |  ![Illustration For loop resource growth with N](/For-count.gif)


## Running Tests
To run tests on different algorithms and datastructures, replace the section `Your function to be tested here` with the algorithm or datastructure (I reccomend encasing it in a function to make commenting out the tested function easier). 

```go
// Element depth to test
for d := 0; d < 6; d++ {
    // Run the function with n elements
    n := int(10 * math.Pow(10, float64(d)))

    // Number of times to test any given element depth
    for test_i := 0; test_i < 10; test_i++ {
        startTime := time.Now()


        /*
            Your function to be tested here
            Your function to be tested here
            Your function to be tested here
        */


        elapsedTime := time.Since(startTime)

        // Call the function to write mem stats to CSV
        writeMemStatsToCSV(writeCounter, n, elapsedTime.Microseconds())

        // Displays the results
        fmt.Printf("Current runtime: %dmicroseconds\nElement depth: %d\n\n", elapsedTime.Microseconds(), n)

        // Runs garbage collection to clear mem-usage before next run
        runtime.GC()

        // Increments writes after test is finished
        writeCounter++
    }
}
```

## Visualization
After processing the data, the application renders a 3D bar chart (`bar3d.html`), where each axis represents a different dimension of the data:

`X-axis`: Different memory metrics.

`Y-axis`: Test sequences of rising N elements.

`Z-axis`: The value of each metric.


## Data Structure
The CSV file `memstats.csv` should have the following columns, representing various memory metrics across different tests:
```
N,Time,Alloc,TotalAlloc,Sys,Lookups,Mallocs,Frees,
HeapAlloc,HeapSys,HeapIdle,HeapInuse,HeapReleased,
HeapObjects,StackInuse,StackSys,MSpanInuse,MSpanSys,
MCacheInuse,MCacheSys,BuckHashSys,GCSys,OtherSys,NextGC,LastGC
```