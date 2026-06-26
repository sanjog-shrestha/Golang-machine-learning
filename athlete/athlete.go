// Package athlete defines the core data types and behaviors.
package athlete

import "fmt"

// Athlete is the "base" struct. Other types embed it to inherit its fields.
type Athlete struct {
	Name   string `json:"name"`
	Team   string `json:"team"`
	Active bool   `json:"active"`
}

// Describe is a method with a value receiver. Embedding types get this for free.
func (a Athlete) Describe() string {
	status := "retired"
	if a.Active {
		status = "active"
	}
	return fmt.Sprintf("%s plays for %s (%s)", a.Name, a.Team, status)
}

// Performer is an INTERFACE. Any type with a Stats() method satisfies it.
// This is how Go does polymorphism — no class hierarchy required.
type Performer interface {
	Stats() string
}
