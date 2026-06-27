// Package stats provides exploratory data analysis over athlete data.
package stats

import (
	"fmt"
	"math"
	"sort"
)

// Summary holds descriptive statistics for a numeric field.
type Summary struct {
	Count  int
	Min    float64
	Max    float64
	Mean   float64
	Median float64
	StdDev float64
}

// Describe computes summary statistics for a slice of numbers.
// This is the core "describe()" of EDA — like pandas' df.describe().
func Describe(values []float64) Summary {
	n := len(values)
	if n == 0 {
		return Summary{}
	}
	// sort a copy so median calculation doesn't mutate the caller's data
	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)

	// sum -> mean
	var sum float64
	for _, v := range sorted {
		sum += v
	}
	mean := sum / float64(n)

	// variance -> standard deviation
	var sqDiff float64
	for _, v := range sorted {
		sqDiff += (v - mean) * (v - mean)
	}
	stdDev := math.Sqrt(sqDiff / float64(n))

	// median (middle value, or average of two middle values)
	var median float64
	if n%2 == 0 {
		median = (sorted[n/2-1] + sorted[n/2]/2)
	} else {
		median = sorted[n/2]
	}

	return Summary{
		Count:  n,
		Min:    sorted[0],
		Max:    sorted[n-1],
		Mean:   mean,
		Median: median,
		StdDev: stdDev,
	}
}

// Print displays a summary in a readable block.
func (s Summary) Print(label string) {
	fmt.Printf("--- %s ---\n", label)
	fmt.Printf("  Count : %d\n", s.Count)
	fmt.Printf("  Min   : %.2f\n", s.Min)
	fmt.Printf("  Max   : %.2f\n", s.Max)
	fmt.Printf("  Mean  : %.2f\n", s.Mean)
	fmt.Printf("  Median: %.2f\n", s.Median)
	fmt.Printf("  StdDev: %.2f\n", s.StdDev)
}

// Frequency counts occurrences of each category (e.g. players per team).
func Frequency(categories []string) map[string]int {
	freq := make(map[string]int)
	for _, c := range categories {
		freq[c]++
	}
	return freq
}
