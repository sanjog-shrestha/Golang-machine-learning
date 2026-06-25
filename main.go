package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Struct tags map JSON keys to GO fields
type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:age`
	Email string `json:email`
}

// clean normalizes a single record in place.
func (p *Person) clean() {
	// Trim whitespace and normalize casing
	p.Name = strings.TrimSpace(p.Name)
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))

	// Capitalize first letter of name if present
	if p.Name != "" {
		p.Name = strings.ToUpper(p.Name[:1]) + p.Name[1:]
	}
}

// valid reports whether a record passes our quality rules.
func (p Person) valid() bool {
	if p.Name == "" {
		return false
	}
	if p.Age < 0 || p.Age > 120 {
		return false
	}
	if !strings.Contains(p.Email, "@") {
		return false
	}
	return true
}

// --- EDA: Summary Statistics over a numeric column ---

type Stats struct {
	Count  int
	Min    float64
	Max    float64
	Mean   float64
	Median float64
	StdDev float64
}

func describe(values []float64) Stats {
	n := len(values)
	if n == 0 {
		return Stats{}
	}

	// Sort a copy for min, max, median
	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)

	var sum float64
	for _, v := range sorted {
		sum += v
	}
	mean := sum / float64(n)

	// Median
	var median float64
	if n%2 == 0 {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	} else {
		median = sorted[n/2]
	}

	// Standard deviation (population)
	var sqDiff float64
	for _, v := range sorted {
		sqDiff += (v - mean) * (v - mean)
	}
	stdDev := math.Sqrt(sqDiff / float64(n))

	return Stats{
		Count:  n,
		Min:    sorted[0],
		Max:    sorted[n-1],
		Mean:   mean,
		Median: median,
		StdDev: stdDev,
	}
}

func (s Stats) print(label string) {
	fmt.Printf("\n--- EDA: %s ---\n", label)
	fmt.Printf("Count:		%d\n", s.Count)
	fmt.Printf("Min:		%.2f\n", s.Min)
	fmt.Printf("Max:		%.2f\n", s.Max)
	fmt.Printf("Mean:		%.2f\n", s.Mean)
	fmt.Printf("Median:		%.2f\n", s.Median)
	fmt.Printf("Std Dev:	%.2f\n", s.StdDev)
}

// -- Visualization -- render an age bar chart to HTML ---
func renderChart(people []Person) error {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Age by Person",
			Subtitle: "Generated with go-echarts",
		}),
	)

	names := make([]string, len(people))
	ages := make([]opts.BarData, len(people))
	for i, p := range people {
		names[i] = p.Name
		ages[i] = opts.BarData{Value: p.Age}
	}

	bar.SetXAxis(names).AddSeries("Age", ages)

	f, err := os.Create("chart.html")
	if err != nil {
		return err
	}

	defer f.Close()
	return bar.Render(f)
}

func main() {
	// Read the whole file into memory
	data, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("read error:", err)
		os.Exit(1)
	}

	// Unmarshal JSON bytes into a slice of structs
	var raw []Person
	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d raw records\n", len(raw))

	// -- Cleansing pipeline ---
	seen := make(map[int]bool) // track IDs to drop duplicates
	var clean []Person
	var dropped int

	// Iterate and process
	for _, p := range raw {
		p.clean() // normalize first

		if !p.valid() || seen[p.ID] {
			dropped++
			continue
		}
		seen[p.ID] = true
		clean = append(clean, p)
	}

	fmt.Printf("Keep %d clean records, dropped %d\n\n", len(clean), dropped)

	// --- EDA ---
	ages := make([]float64, len(clean))
	for i, p := range clean {
		ages[i] = float64(p.Age)
	}
	describe(ages).print("Age")

	// --- Visualization ---
	if err := renderChart(clean); err != nil {
		fmt.Println("chart error:", err)
		os.Exit(1)
	}
	fmt.Println("\nWrote chart.html")
}
