package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go-sports/athlete"
	"go-sports/preprocess"
	"go-sports/stats" // NEW
	"go-sports/viz"   // NEW
)

type Roster struct {
	Footballers []athlete.Footballer `json:"footballers"`
	Cricketers  []athlete.Cricketer  `json:"cricketers"`
}

func main() {
	data, err := os.ReadFile("players.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var roster Roster
	if err := json.Unmarshal(data, &roster); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// --- Cleansing (from previous step) ---
	roster.Footballers = preprocess.CleanFootballers(roster.Footballers)
	roster.Cricketers = preprocess.CleanCricketers(roster.Cricketers)

	// ============ EXPLORATORY DATA ANALYSIS ============
	fmt.Println("===== EDA: Football Goals =====")

	// extract the numeric field we want to analyze
	goals := make([]float64, 0, len(roster.Footballers))
	for _, f := range roster.Footballers {
		goals = append(goals, float64(f.Goals))
	}
	goalSummary := stats.Describe(goals)
	goalSummary.Print("Goals")

	// categorical EDA: how many players per team
	teams := make([]string, 0)
	for _, f := range roster.Footballers {
		teams = append(teams, f.Team)
	}
	freq := stats.Frequency(teams)
	fmt.Println("\n--- Players per Team ---")
	for team, n := range freq {
		fmt.Printf("  %s: %d\n", team, n)
	}

	// ============ DATA VISUALIZATION ============
	// Build a bar chart of goals per footballer.
	bars := make([]viz.Bar, 0, len(roster.Footballers))
	for _, f := range roster.Footballers {
		bars = append(bars, viz.Bar{Label: f.Name, Value: float64(f.Goals)})
	}

	html := viz.BarChartHTML("Goals by Footballer", bars)
	if err := viz.WriteHTML("goals_chart.html", html); err != nil {
		fmt.Println("Error writing chart:", err)
		return
	}
	fmt.Println("\nChart written to goals_chart.html (open it in a browser)")
}
