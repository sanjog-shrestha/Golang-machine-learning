// Package preprocess cleans and normalizes raw athlete data
// before it is used by the rest of the application.
package preprocess

import (
	"go-sports/athlete"
	"regexp"
	"strings"
)

// collapseSpaces matches one or more whitespace characters.
var collapseSpaces = regexp.MustCompile(`\s+`)

// cleanString trims, collapses internal whitespace, and title-cases text.
// "  lionel   MESSI " -> "Lionel Messi"
func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = collapseSpaces.ReplaceAllString(s, " ")
	return strings.Title(strings.ToLower(s))
}

// CleanFootballers normalizes fields, drops invalid rows, and removes duplicates.
func CleanFootballers(in []athlete.Footballer) []athlete.Footballer {
	seen := make(map[string]bool) // tracks unique name+team keys
	var out []athlete.Footballer

	for _, f := range in {
		// 1. NORMALIZE textual fields
		f.Name = cleanString(f.Name)
		f.Team = cleanString(f.Team)
		f.Sport = cleanString(f.Sport)

		// 2. VALIDATE: drop rows with no name (missing required data)
		if f.Name == "" {
			continue
		}

		// 3. SANITIZE numeric fields: a negative goal count is invalid -> clamp to 0
		if f.Goals < 0 {
			f.Goals = 0
		}

		// 4. DEDUPLICATE on a composite key
		key := f.Name + "|" + f.Team
		if seen[key] {
			continue
		}

		seen[key] = true

		out = append(out, f)
	}

	return out

}

// CleanCricketers applies the same pipeline to cricketers.
func CleanCricketers(in []athlete.Cricketer) []athlete.Cricketer {
	seen := make(map[string]bool)
	var out []athlete.Cricketer

	for _, c := range in {
		c.Name = cleanString(c.Name)
		c.Team = cleanString(c.Team)
		c.Sport = cleanString(c.Sport)

		if c.Name == "" {
			continue
		}

		if c.Runs < 0 {
			c.Runs = 0
		}

		// 4. DEDUPLICATE on a composite key
		key := c.Name + "|" + c.Team
		if seen[key] {
			continue
		}

		seen[key] = true

		out = append(out, c)
	}

	return out

}
