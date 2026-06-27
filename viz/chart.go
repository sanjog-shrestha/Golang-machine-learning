// Package viz renders simple SVG bar charts embedded in an HTML file.
package viz

import (
	"fmt"
	"os"
	"strings"
)

// Bar is a single labeled value in a chart.
type Bar struct {
	Label string
	Value float64
}

// BarChartHTML builds a self-contained HTML page with an inline SVG bar chart.
func BarChartHTML(title string, bars []Bar) string {
	const (
		width    = 600
		barH     = 36
		gap      = 12
		leftPad  = 140 // room for labels
		rightPad = 60  // room for value text
	)
	height := len(bars)*(barH+gap) + gap

	// find the max value to scale bars proportionally
	var max float64
	for _, b := range bars {
		if b.Value > max {
			max = b.Value
		}
	}
	if max == 0 {
		max = 1 // avoid divide-by-zero
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" font-family="sans-serif">`,
		width, height))

	for i, b := range bars {
		y := gap + i*(barH+gap)
		barW := (b.Value / max) * float64(width-leftPad-rightPad)

		// label (left), bar (middle), value (right)
		sb.WriteString(fmt.Sprintf(
			`<text x="%d" y="%d" font-size="14" text-anchor="end" dominant-baseline="middle">%s</text>`,
			leftPad-10, y+barH/2, b.Label))
		sb.WriteString(fmt.Sprintf(
			`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="#3b82f6"/>`,
			leftPad, y, barW, barH))
		sb.WriteString(fmt.Sprintf(
			`<text x="%.1f" y="%d" font-size="13" dominant-baseline="middle">%.0f</text>`,
			float64(leftPad)+barW+8, y+barH/2, b.Value))
	}
	sb.WriteString(`</svg>`)

	return fmt.Sprintf(
		`<!DOCTYPE html><html><head><meta charset="utf-8"><title>%s</title></head>`+
			`<body><h2 style="font-family:sans-serif">%s</h2>%s</body></html>`,
		title, title, sb.String())
}

// WriteHTML saves the chart HTML to a file.
func WriteHTML(path, html string) error {
	return os.WriteFile(path, []byte(html), 0644)
}
