# go-sports

A small, beginner-friendly Go project for learning the full data workflow — **loading, cleansing, exploratory analysis, visualization, and a from-scratch machine learning model** — with a sports theme, packaged to run in Docker. It demonstrates Go modules, struct embedding (Go's take on inheritance), interfaces, a preprocessing pipeline, descriptive statistics, dependency-free SVG charting, and a linear regression model built without any ML libraries.

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Quick Start (Docker)](#quick-start-docker)
- [Running Locally (without Docker)](#running-locally-without-docker)
- [How It Works](#how-it-works)
- [Packages](#packages)
- [The Linear Regression Model](#the-linear-regression-model)
- [Go Concepts Demonstrated](#go-concepts-demonstrated)
- [Sample Output](#sample-output)
- [Next Steps](#next-steps)

---

## Features

- Load sports data from a JSON file into typed Go structs.
- Organize code into **modules and packages** (`athlete`, `preprocess`, `stats`, `viz`, `mlmodel`).
- Model data with **struct embedding** (composition instead of inheritance).
- Use **interfaces** for polymorphism across different athlete types.
- A **data cleansing & preprocessing** pipeline that normalizes, validates, sanitizes, and deduplicates messy input.
- **Exploratory Data Analysis (EDA)**: descriptive statistics (count, min, max, mean, median, std dev) plus categorical frequency counts.
- **Data visualization**: a self-contained SVG bar chart written to an HTML file — no external dependencies.
- **Linear regression from scratch**: hypothesis, cost function, gradient descent, training, prediction, and evaluation — no ML libraries.
- A **multi-stage Docker build** producing a small Alpine-based runtime image.

---

## Project Structure

```
go-sports/
├── Dockerfile            # Multi-stage build (compile in Go image, run on Alpine)
├── go.mod                # Module definition: module path "go-sports"
├── main.go               # Entry point: load → clean → analyze → visualize → train
├── players.json          # Sample (deliberately messy) sports data
├── athlete/              # Core domain package
│   ├── athlete.go        # Base Athlete struct + Performer interface
│   └── specialized.go    # Footballer & Cricketer (embed Athlete)
├── preprocess/           # Data cleansing package
│   └── clean.go          # Normalize, validate, sanitize, deduplicate
├── stats/                # Exploratory Data Analysis package
│   └── eda.go            # Describe() summary stats + Frequency() counts
├── viz/                  # Visualization package
│   └── chart.go          # SVG bar chart embedded in standalone HTML
└── mlmodel/              # Machine learning package
    └── regression.go     # Linear regression from scratch
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

This prints the EDA summary and regression results to the console and drops `goals_chart.html` right in your project folder. Or build a binary:

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
4. **Analyze (EDA)** — the `stats` package computes descriptive statistics for numeric fields and frequency counts for categories.
5. **Visualize** — the `viz` package builds an SVG bar chart and writes it to `goals_chart.html`.
6. **Train & predict** — the `mlmodel` package fits a linear regression model with gradient descent, evaluates it, and makes a prediction.

The sample `players.json` is intentionally messy — stray whitespace, mixed casing, a negative goal count, an empty name, and duplicate players — so the cleansing step has real work to do before analysis.

---

## Packages

| Package | Responsibility | Key types / functions |
|---|---|---|
| `athlete` | Domain model | `Athlete`, `Footballer`, `Cricketer`, `Performer` interface |
| `preprocess` | Data cleansing | `CleanFootballers`, `CleanCricketers` |
| `stats` | Exploratory analysis | `Describe` → `Summary`, `Frequency` |
| `viz` | Visualization | `Bar`, `BarChartHTML`, `WriteHTML` |
| `mlmodel` | Machine learning | `Sample`, `LinearModel`, `Normalize` |

### `stats` — Exploratory Data Analysis

`Describe([]float64)` returns a `Summary` with count, min, max, mean, median, and standard deviation — the rough equivalent of pandas' `df.describe()`, built from scratch using the standard `math` and `sort` packages. It copies the input before sorting so the caller's data is never mutated. `Frequency([]string)` returns a `map[string]int` of category counts.

### `viz` — Data Visualization

`BarChartHTML(title, bars)` assembles an inline SVG bar chart (using `strings.Builder`) and wraps it in a standalone HTML page. Bars are scaled proportionally to the largest value, with a divide-by-zero guard. `WriteHTML(path, html)` saves it to disk. No charting library is required.

### `mlmodel` — Linear Regression

A complete simple linear regression (`y = w*x + b`) implemented from scratch, predicting goals from matches played. See the dedicated section below.

---

## The Linear Regression Model

The `mlmodel` package implements all six stages of building a regression model, with no external ML dependencies.

| Stage | Implementation | What it does |
|---|---|---|
| **1. Model & data** | `Sample`, `LinearModel` | A `Sample` is one (feature, target) pair; `LinearModel` holds `Weight` (slope) and `Bias` (intercept) |
| **2. Hypothesis** | `Predict(x)` | Computes `h(x) = w*x + b`, the model's prediction |
| **3. Cost function** | `Cost(data)` | Mean Squared Error: `J = (1/2n) Σ (h(xᵢ) − yᵢ)²` |
| **4. Gradient descent** | `gradientStep(data, lr)` | Computes partial derivatives and nudges parameters downhill |
| **5. Training** | `Train(data, lr, epochs)` | Repeats gradient steps, returns the cost history |
| **6. Prediction & evaluation** | `Predict`, `RSquared`, `RMSE`, `Summary` | Forecasts new values and reports goodness of fit |

### Key ideas

- **Mean Squared Error** averages the squared gap between predictions and truth. Squaring penalizes large errors and stays positive; the `1/2` factor cancels cleanly when differentiating.
- **Gradient descent** moves parameters opposite to the cost gradient. The **learning rate** sets step size — too large diverges, too small crawls. Training runs for a number of **epochs** and records cost so you can confirm convergence.
- **Evaluation** uses **R²** (variance explained; 1.0 is perfect, 0 is no better than the mean) and **RMSE** (average error in the original units).
- **Feature scaling matters.** `Normalize` applies min-max scaling so gradient descent converges reliably; unscaled features are the most common cause of a from-scratch model diverging. To predict on a new raw input, scale it the same way before calling `Predict`.
- **Pointer vs. value receivers.** Read-only methods (`Predict`, `Cost`, `RSquared`) use value receivers; the mutating `gradientStep` uses a pointer receiver `(m *LinearModel)` so it can update the parameters in place.

---

## Go Concepts Demonstrated

| Concept | Where | What to look for |
|---|---|---|
| Modules & packages | `go.mod`, all folders | Module path `go-sports`; subpackages imported as `go-sports/mlmodel` |
| Structs & struct tags | `athlete/*.go` | `` `json:"name"` `` tags map struct fields to JSON keys |
| Embedding (composition) | `specialized.go` | `Footballer` embeds `Athlete` and inherits its fields/methods |
| Field & method promotion | `main.go` | `footballer.Name` and `footballer.Describe()` come from `Athlete` |
| Interfaces (polymorphism) | `athlete.go` | `Performer` interface; satisfied implicitly by any type with `Stats()` |
| Methods & receivers | `athlete.go`, `eda.go`, `regression.go` | Value receivers for reads; pointer receiver for `gradientStep` |
| JSON marshal/unmarshal | `main.go` | `encoding/json` for decoding |
| Error handling | `main.go` | The idiomatic `if err != nil` pattern |
| Maps as sets / counters | `clean.go`, `eda.go` | `map[string]bool` to dedupe; `map[string]int` to count |
| Regular expressions | `preprocess/clean.go` | `regexp` collapses repeated whitespace |
| Math & sorting | `stats/eda.go`, `regression.go` | `math.Sqrt`, `sort.Float64s` |
| Efficient string building | `viz/chart.go` | `strings.Builder` to assemble SVG |
| Multiple return values | `mlmodel/regression.go` | `Normalize` returns scaled data plus min and span |
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

===== Linear Regression: predict goals from matches =====
Cost: start=0.1234  end=0.0012
Learned model: goals = 0.470 * matches + 0.030
  R²   : 0.998
  RMSE : 0.350

Prediction: a player with 35 matches scores ~16.1 goals
```

Exact numbers vary, but cost should fall sharply from start to end and R² should be close to 1.0 on this clean, near-linear sample.

---

## Next Steps

Ideas to extend the project as you keep learning:

- Add a **cost-history chart** to `viz` so you can see the model converge.
- Extend to **multiple features** (multivariate regression).
- Add a **train/test split** to evaluate on unseen data.
- Add a **histogram** to show the distribution of goals.
- Make cleansing **configurable** with a `CleanOptions` struct.
- Read from **CSV** instead of JSON to practice another format.
- Expose the cleaned data, charts, and predictions over an **HTTP API**.
- Add **unit tests** (`go test`) for the cleansing, stats, and model functions.

---

## License

This is a learning project — use it freely.
