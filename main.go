package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Struct tags map JSON keys to GO fields
type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:age`
}

func main() {
	// Read the whole file into memory
	data, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("read error:", err)
		os.Exit(1)
	}

	// Unmarshal JSON bytes into a slice of structs
	var people []Person
	if err := json.Unmarshal(data, &people); err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	// Iterate and process
	for _, p := range people {
		fmt.Printf("%d: %s (%d)\n", p.ID, p.Name, p.Age)
	}

	// Filter example: people over 28
	var older []Person
	for _, p := range people {
		if p.Age > 28 {
			older = append(older, p)
		}
	}

	fmt.Printf("\nOver 28: %d people\n", len(older))

	// Marshal back to JSON (writing data out)
	out, _ := json.MarshalIndent(older, "", " ")
	fmt.Println(string(out))
}
