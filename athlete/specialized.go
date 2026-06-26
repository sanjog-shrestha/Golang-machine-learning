package athlete

import "fmt"

// Footballer EMBEDS Athlete. This is composition, Go's alternative to inheritance.
// A Footballer automatically gets Name, Team, Active, and the Describe() method.
type Footballer struct {
	Athlete        // embedded — no field name, this is the key part
	Goals   int    `json:"goals"`
	Sport   string `json:"sport`
}

// Stats satisfies the Performer interface for Footballer.
func (f Footballer) Stats() string {
	return fmt.Sprintf("%d goals", f.Goals)
}

// Cricketer also embeds Athlete but reports different stats.
type Cricketer struct {
	Athlete        // embedded — no field name, this is the key part
	Runs    int    `json:"runs"`
	Sport   string `json:"sport`
}

func (c Cricketer) Stats() string {
	return fmt.Sprintf("%d runs", c.Runs)
}
