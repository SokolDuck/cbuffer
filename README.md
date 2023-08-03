# CircuitBuffer Go Module

This Go module provides two data structures: **`CircuitBuffer`** and **`OrderedCircuitBuffer`**. Both are circular buffers, or ring buffers, which are data structures that use a single, fixed-size buffer as if it were connected end-to-end. The **`OrderedCircuitBuffer`** maintains the order of elements and can perform a binary search.

These buffers can store any type that implements the **`Comparable`** interface.

```golang
type Comparable[T any] interface {
	Less(T) bool
	Equal(T) bool
}

type CircuitBuffer[T Comparable[T]];
type OrderedCircuitBuffer[T Comparable[T]];
```

## Features
- When the buffer is full and a new item is added, the oldest item will be removed from the buffer.

- Size of buffer is fixed, this mean we never will be allocate and copy extra memory.

## Main Functions

- `NewCircuitBuffer[T Comparable[T]](size int) CircuitBuffer[T}]`: Creates a new CircuitBuffer of a given size.

- `NewOrderedCircuitBuffer[T Comparable[T]](size int) OrderedCircuitBuffer[T]`: Creates a new OrderedCircuitBuffer of a given size.

### CircuitBuffer Functions 
- `Add(item T) bool`: Adds a new item to the buffer. If the buffer is full, it will replace the oldest item.

- `GetItem(index int) T`: Return item by index.

- `Len() int`: Returns the current length of the payload.

- `Cap() int`: Returns the capacity of the buffer.

- `Iter() chan *T`: Returns a channel to iterate over the items in the buffer. Chan **return POINTER** on item.

- `Break()`: Breaks the iteration early and cleans up the iterator.

- `String() string`: To support "%s" and conversion to string.

### OrderedCircuitBuffer Functions
- All the same as CircuitBuffer

- `Add(item T) error`: Adds a new item to the buffer. The item must be greater than or equal to the last item in the buffer.

- `Search(value T) index int, found bool`: Performs a binary search for a value in the buffer. Returns the index of the value and a boolean indicating whether the value was found.

> **`All functions are thread safe.`**

## Limitations

- The binary search in **`OrderedCircuitBuffer`** assumes that the buffer is sorted in ascending order. If this condition is not met, the search results will be incorrect.

- Attempting to add an item to an **`OrderedCircuitBuffer`** that is less than the last item in the buffer will return an error.

- The capacity of the buffer is set at creation and cannot be changed afterwards.


## Use Cases
This data structure can be useful for any task related with sliding window.

Here are a few examples:

1. **Log Buffering**

	Imagine you have an application that constantly logs data, but you only care about preserving the last N log entries due to memory constraints. You could use the CircuitBuffer to keep only the most recent N log entries in memory, dropping older entries as necessary.

2. **Time Series Data**

	If you're working with time series data where data points arrive in an ordered fashion (e.g. stock prices, sensor readings), you might want to keep a buffer of the latest N points for quick calculations like moving averages. Here, OrderedCircuitBuffer would be appropriate. You could also search for specific data points in the buffer using binary search, which is faster than linear search.

3. **Streaming Data**

	If you're dealing with a continuous stream of data and need to take action on each piece of data as it comes in, you could use the Iter() method provided by this module to process each data point in the buffer. This can be useful in a variety of real-time data processing scenarios.

4. **Producer-Consumer Problem**

	The CircuitBuffer can be used to solve the producer-consumer problem in concurrent programming. The producer can continuously add items to the buffer, and the consumer can read items from the buffer. If the buffer is full, the producer can overwrite the oldest items.


## Examples
### 1. Log Buffering

```go
package main

import (
	"fmt"

	cb "github.com/SokolDuck/cbuffer"
)

func ProcessLogs(buffer *cb.CircuitBuffer[string]) {
	// Call buffer.Break() if the Iter loop stops before it runs out of data. Or call it with "deffer" to make sure you don't spawn dead goroutines.
	defer buffer.Break()

	for v := range buffer.Iter() {
		fmt.Println(*v)
	}

}

func main() {
	cb := cb.NewCircuitBuffer[string](5)

	logs := []string{"Log1", "Log2", "Log3", "Log4", "Log5", "Log6", "Log7"}

	for _, log := range logs {
		cb.Add(log)
	}

	ProcessLogs(cb)
}
```

### 2. Time Series Data

```go
package main

import (
	"fmt"

	cb "github.com/SokolDuck/cbuffer"
)

type TimeSeriesData struct {
	timestamp int
	data      int
}

func (ts TimeSeriesData) Less(other TimeSeriesData) bool {
	return ts.timestamp < other.timestamp
}

func (ts TimeSeriesData) Equal(other TimeSeriesData) bool {
	return ts.timestamp == other.timestamp
}

func main() {
	cbuf := cb.NewOrderedCircuitBuffer[TimeSeriesData](5)

	dataPoints := make([]TimeSeriesData, 8)

	for i := 0; i < 8; i++ {
		tsd := TimeSeriesData{
			timestamp: i,
			data:      i * i,
		}
		dataPoints = append(dataPoints, tsd)
	}

	for _, dataPoint := range dataPoints {
		err := cbuf.Add(dataPoint)
		if err != nil {
			fmt.Println("Error occurred while adding:", err)
		}
	}

	for v := range cbuf.Iter() {
		fmt.Println(*v)

		if v.data > 10 {
			cbuf.Break()
			break // or continue - behavior will be the same after calling cb.Break()
		}
	}

	searchItem := new(TimeSeriesData)
	searchItem.timestamp = 4

	index, found := cbuf.Search(*searchItem)
	fmt.Printf("buf: %v\n", cbuf)
	fmt.Printf("Searching for 4, Index: %v, Found: %v\n", index, found)
}

```
