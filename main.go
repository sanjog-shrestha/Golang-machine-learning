package main

import (
	"encoding/json"
	"fmt"
	"go-sports/athlete"
	"go-sports/preprocess"
	"os"
)

// Roster matches the top-level shape of the JSON file.
type Roster struct {
	Footballers []athlete.Footballer `json:"footballers"`
	Cricketers  []athlete.Cricketer  `json:"Cricketers"`
}

func main() {
	data, err := os.ReadFile("players.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var roster Roster
	if err := json.Unmarshal(data, &roster); err != nil {
		fmt.Println("Error parsing json:", err)
		return
	}

	// --- RAW counts before cleansing ---
	fmt.Printf("Raw: %d footballers. %d cricketers\n", len(roster.Footballers), len(roster.Cricketers))

	// --- CLEANSING + PREPROCESSING STEP ---
	roster.Footballers = preprocess.CleanFootballers(roster.Footballers)
	roster.Cricketers = preprocess.CleanCricketers(roster.Cricketers)

	fmt.Printf("Cleaned: %d footballers, %d cricketers\n", len(roster.Footballers), len(roster.Cricketers))

	// Collect everyone into a slice of the Performer INTERFACE.
	// Different concrete types living in one collection = polymorphism.
	var performers []athlete.Performer
	for _, f := range roster.Footballers {
		performers = append(performers, f)
	}

	for _, c := range roster.Cricketers {
		performers = append(performers, c)
	}

	fmt.Println("=== Ckeaned Perfomers ===")
	for _, f := range roster.Footballers {
		fmt.Printf("%s | %s\n", f.Describe(), f.Stats())
	}

	for _, c := range roster.Cricketers {
		fmt.Printf("%s | %s\n", c.Describe(), c.Stats())
	}

}
