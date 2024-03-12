import (
    "fmt"
    "runtime"
)

func main() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("Stack usage: %v bytes\n", m.TotalAlloc - m.HeapAlloc)

var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("HeapAlloc: %v bytes\n", m.HeapAlloc)
}

