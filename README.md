# gostatistics

`gostatistics` is a lightweight Go library that provides common descriptive statistics with a friendly API. It is designed for quick analyses, report generation, and teaching scenarios where you just need reliable numbers without pulling in a heavy dependency.

## Installation

```bash
go get github.com/evdnx/gostatistics
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/evdnx/gostatistics"
)

func main() {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}

	stats, err := gostatistics.Describe(values)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Count: %d\n", stats.Count)
	fmt.Printf("Mean: %.2f\n", stats.Mean)
	fmt.Printf("Median: %.2f\n", stats.Median)
	fmt.Printf("Sample Std Dev: %.2f\n", stats.SampleStandardDeviation)
	fmt.Printf("Range: %.2f\n", stats.Range)
}
```

Output:

```
Count: 8
Mean: 5.00
Median: 4.50
Sample Std Dev: 2.14
Range: 7.00
```

## Features

- Arithmetic helpers: `Sum`, `Mean`, `Range`, `Min`, `Max`, `Median`
- Variance and standard deviation (sample and population)
- Convenience wrappers: `SampleVariance`, `PopulationVariance`, `SampleStandardDeviation`, `PopulationStandardDeviation`
- Dataset summary via `Describe` with ready-to-print descriptive statistics
- Configurable behaviour through `StatisticsOptions`
- Exported errors (`ErrEmptySlice`, `ErrInsufficientValues`) to make error handling explicit

## Working with options

Most variance-based functions accept an optional `StatisticsOptions` value. By default the library uses Bessel’s correction (`n-1`) so you get sample statistics out of the box. To compute population values you can pass:

```go
opts := gostatistics.StatisticsOptions{
	UseSampleCorrection: false,
}
variance, err := gostatistics.Variance(values, opts)
```

Alternatively, call one of the convenience helpers (`PopulationVariance`, `PopulationStandardDeviation`), which set the option for you.

## Error handling

- `ErrEmptySlice` — returned when an operation requires at least one value.
- `ErrInsufficientValues` — returned by variance/standard deviation calculations that need at least two values.

Use [`errors.Is`](https://pkg.go.dev/errors#Is) to branch on these sentinel errors.

## Testing

Run the full suite with:

```bash
go test ./...
```
