# go-sports

A small, beginner-friendly Go project for learning **data loading, handling, cleansing, exploratory analysis, and visualization** with a sports theme — packaged to run in Docker. It demonstrates Go modules, struct embedding (Go's take on inheritance), interfaces, a data preprocessing pipeline, descriptive statistics, and dependency-free SVG charting.

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Quick Start (Docker)](#quick-start-docker)
- [Running Locally (without Docker)](#running-locally-without-docker)
- [How It Works](#how-it-works)
- [Packages](#packages)
- [Go Concepts Demonstrated](#go-concepts-demonstrated)
- [Sample Output](#sample-output)
- [Next Steps](#next-steps)

---

## Features

- Load sports data from a JSON file into typed Go structs.
- Organize code into **modules and packages** (`athlete`, `preprocess`, `stats`, `viz`).
- Model data with **struct embedding** (composition instead of inheritance).
- Use **interfaces** for polymorphism across different athlete types.
- A **data cleansing & preprocessing** pipeline that normalizes, validates, sanitizes, and deduplicates messy input.
- **Exploratory Data Analysis (EDA)**: descriptive statistics (count, min, max, mean, median, std dev) plus categorical frequency counts.
- **Data visualization**: a self-contained SVG bar chart written to an HTML file — no external dependencies.
- A **multi-stage Docker build** producing a small Alpine-based runtime image.

---

## Project Structure

```
go-sports/
├── Dockerfile            # Multi-stage build (compile in Go image, run on Alpine)
├── go.mod                # Module definition: module path "go-sports"
├── main.go               # Entry point: load → clean → analyze → visualize
├── players.json          # Sample (deliberately messy) sports data
├── athlete/              # Core domain package
│   ├── athlete.go        # Base Athlete struct + Performer interface
│   └── specialized.go    # Footballer & Cricketer (embed Athlete)
├── preprocess/           # Data cleansing package
│   └── clean.go          # Normalize, validate, sanitize, deduplicate
├── stats/                # Exploratory Data Analysis package
│   └── eda.go            # Describe() summary stats + Frequency() counts
└── viz/                  # Visualization package
    └── chart.go          # SVG bar chart embedded in standalone HTML
```

---

## Prerequisites

Choose one path:

- **Docker** — the only requirement for the recommended workflow. [Install Docker](https://docs.docker.com/get-docker/).
- **Go 1.22+** — if you prefer to run it directly on your machine. [Install Go](https://go.dev/dl/).

---

## Quick Start (Docker)

Build the image:

```bash
docker build -t go-sports .
```

Run the container. Mount the current folder so the generated chart lands on your machine:

```bash
docker run --rm -v "$(pwd)":/app/out go-sports sh -c "./sports-app && cp goals_chart.html /app/out/"
```

Then open `goals_chart.html` in a browser.

---

## Running Locally (without Docker)

From the project root:

```bash
go run .
```

This prints the EDA summary to the console and drops `goals_chart.html` right in your project folder. Or build a binary:

```bash
go build -o sports-app .
./sports-app
```

---

## How It Works

The program runs a pipeline in `main.go`:

1. **Load** — `os.ReadFile` reads `players.json` as raw bytes.
2. **Parse** — `json.Unmarshal` decodes the JSON into a `Roster` struct.
3. **Cleanse** — the `preprocess` package normalizes text, drops invalid rows, sanitizes numbers, and removes duplicates.
4. **Analyze (EDA)** — the `stats` package computes descriptive statistics for numeric fields (e.g. goals) and frequency counts for categories (e.g. players per team).
5. **Visualize** — the `viz` package builds an SVG bar chart and writes it to `goals_chart.html`.

The sample `players.json` is intentionally messy — stray whitespace, mixed casing, a negative goal count, an empty name, and duplicate players — so the cleansing step has real work to do before analysis.

---

## Packages

| Package | Responsibility | Key types / functions |
|---|---|---|
| `athlete` | Domain model | `Athlete`, `Footballer`, `Cricketer`, `Performer` interface |
| `preprocess` | Data cleansing | `CleanFootballers`, `CleanCricketers` |
| `stats` | Exploratory analysis | `Describe` → `Summary`, `Frequency` |
| `viz` | Visualization | `Bar`, `BarChartHTML`, `WriteHTML` |

### `stats` — Exploratory Data Analysis

`Describe([]float64)` returns a `Summary` with count, min, max, mean, median, and standard deviation — the rough equivalent of pandas' `df.describe()`, built from scratch using the standard `math` and `sort` packages. It copies the input before sorting so the caller's data is never mutated. `Frequency([]string)` returns a `map[string]int` of category counts for categorical analysis.

### `viz` — Data Visualization

`BarChartHTML(title, bars)` assembles an inline SVG bar chart (using `strings.Builder`) and wraps it in a standalone HTML page. Bars are scaled proportionally to the largest value, with a divide-by-zero guard when all values are zero. `WriteHTML(path, html)` saves it to disk. No charting library is required, which keeps the Docker image small.

---

## Go Concepts Demonstrated

| Concept | Where | What to look for |
|---|---|---|
| Modules & packages | `go.mod`, all folders | Module path `go-sports`; subpackages imported as `go-sports/stats` |
| Structs & struct tags | `athlete/*.go` | `` `json:"name"` `` tags map struct fields to JSON keys |
| Embedding (composition) | `specialized.go` | `Footballer` embeds `Athlete` and inherits its fields/methods |
| Field & method promotion | `main.go` | `footballer.Name` and `footballer.Describe()` come from `Athlete` |
| Interfaces (polymorphism) | `athlete.go` | `Performer` interface; satisfied implicitly by any type with `Stats()` |
| Methods & receivers | `athlete.go`, `eda.go` | `func (a Athlete) Describe()`, `func (s Summary) Print()` |
| JSON marshal/unmarshal | `main.go` | `encoding/json` for decoding |
| Error handling | `main.go` | The idiomatic `if err != nil` pattern |
| Maps as sets / counters | `preprocess/clean.go`, `eda.go` | `map[string]bool` to dedupe; `map[string]int` to count |
| Regular expressions | `preprocess/clean.go` | `regexp` collapses repeated whitespace |
| Math & sorting | `stats/eda.go` | `math.Sqrt`, `sort.Float64s` for stats |
| Efficient string building | `viz/chart.go` | `strings.Builder` to assemble SVG |
| Value semantics | `clean.go`, `eda.go` | Slices copied before mutation; functions stay pure |

### A note on "inheritance"

Go has **no classical inheritance**. Instead it uses **composition** via struct embedding (a `Footballer` *has an* `Athlete`) and **interfaces** for shared behavior. This is more explicit than subclassing and avoids fragile deep class hierarchies.

### Data cleansing stages

- **Normalization** — trim/collapse whitespace and unify casing (`MESSI` → `Messi`).
- **Validation** — drop records missing required fields (e.g. empty name).
- **Sanitization** — fix out-of-range values (a negative goal count is clamped to `0`).
- **Deduplication** — collapse repeats using a composite `name|team` key.

---

## Sample Output

Console output follows this shape:

```
===== EDA: Football Goals =====
--- Goals ---
  Count : 2
  Min   : 0.00
  Max   : 15.00
  Mean  : 7.50
  Median: 7.50
  StdDev: 7.50

--- Players per Team ---
  Inter Miami: 1
  Santos: 1

Chart written to goals_chart.html (open it in a browser)
```

Opening `goals_chart.html` shows a horizontal bar chart of goals per footballer, with each bar scaled to the highest scorer.

---

## Next Steps

Ideas to extend the project as you keep learning:

- Add a **histogram** to show the distribution of goals.
- Compute **correlations** between numeric fields as the dataset grows.
- Make cleansing **configurable** with a `CleanOptions` struct.
- Read from **CSV** instead of JSON to practice another format.
- Expose the cleaned data and charts over an **HTTP API**.
- Swap the hand-rolled SVG for **`gonum/plot`** for richer charts.
- Add **unit tests** (`go test`) for the cleansing and stats functions.

---

## License

This is a learning project — use it freely.
