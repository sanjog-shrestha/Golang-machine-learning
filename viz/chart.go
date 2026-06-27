// Package viz renders simple SVG bar charts embedded in an HTML file.
package viz

import (
	"fmt"
	"math"
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

// Point is an (x, y) pair for scatter plots.
type Point struct {
	X, Y float64
}

// ResidualPlotHTML renders residuals as a scatter around a zero baseline.
// Points above the line over-predict; below, under-predict
func ResidualPlotHTML(title string, points []Point) string {
	const (
		w, h = 600, 360
		padL = 50
		padB = 40
		padT = 30
		padR = 20
	)
	plotW := w - padL - padR
	plotH := h - padT - padB

	// find ranges
	minX, maxX := points[0].X, points[0].X
	maxAbsY := 0.0
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if a := math.Abs(p.Y); a > maxAbsY {
			maxAbsY = a
		}
	}
	spanX := maxX - minX
	if spanX == 0 {
		spanX = 1
	}
	if maxAbsY == 0 {
		maxAbsY = 1
	}

	// map data coords -> pixel coords
	sx := func(x float64) float64 { return padL + (x-minX)/spanX*float64(plotW) }
	sy := func(y float64) float64 { return padT + (1-(y+maxAbsY)/(2*maxAbsY))*float64(plotH) }

	zeroY := sy(0)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`<svg width="%d" height="%d" xmlns="http://www.w3.org/200/svg" font-family="sans-serif">`,
		w, h))

	// zero reference line
	sb.WriteString(fmt.Sprintf(
		`<line x1="%d" y1="%.1f" x2="%d" y2="%.1f" stroke="#888" stroke-dasharray="4"/>`,
		padL, zeroY, w-padR, zeroY))

	// points
	for _, p := range points {
		color := "#3b82f6"
		if p.Y < 0 {
			color = "#ef4444"
		}
		sb.WriteString(fmt.Sprintf(
			`<circle cx="%.1f" cy="%.1f" r="5" fill="%s"/>`, sx(p.X), sy(p.Y), color))
	}
	sb.WriteString("</svg>'")

	return fmt.Sprintf(
		`<!DOCTYPE html><html><head><meta charset="utf-8"><title>%s</title></head>`+
			`<body><h2 style="font-family:sans-serif">%s</h2>%s`+
			`<p style="font-family:sans-serif;color:#555">Blue = under-predicted, Red = over-predicted. `+
			`Random scatter around the line = good linear fit.</p></body></html>`,
		title, title, sb.String())
}
