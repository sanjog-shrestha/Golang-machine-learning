# go-sports

A small, beginner-friendly Go project for learning the full data workflow — **loading, cleansing, exploratory analysis, visualization, and a from-scratch machine learning model with fine-tuning** — built around a sports theme and packaged to run in Docker. It demonstrates Go modules, struct embedding (Go's take on inheritance), interfaces, a preprocessing pipeline, descriptive statistics, dependency-free SVG charting, and a linear regression model (with MSE, residual plots, and L1 regularization) built without any ML libraries.

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Quick Start (Docker Compose)](#quick-start-docker-compose)
- [Running with Plain Docker](#running-with-plain-docker)
- [Running Locally (without Docker)](#running-locally-without-docker)
- [How It Works](#how-it-works)
- [Packages](#packages)
- [The Linear Regression Model](#the-linear-regression-model)
- [Fine-Tuning: Evaluation & Optimization](#fine-tuning-evaluation--optimization)
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
- **Data visualization**: self-contained SVG charts (bar chart + residual plot) written to HTML files — no external dependencies.
- **Linear regression from scratch**: hypothesis, cost function, gradient descent, training, prediction, and evaluation — no ML libraries.
- **Model fine-tuning**: standard MSE metric, residual diagnostics, and L1 (Lasso) regularization.
- A **multi-stage Docker build** plus a **Docker Compose** workflow that writes charts back to your host.

---

## Project Structure

```
go-sports/
├── Dockerfile            # Multi-stage build (compile in Go image, run on Alpine)
├── docker-compose.yml    # One-command build/run; mounts ./output for charts
├── go.mod                # Module definition: module path "go-sports"
├── main.go               # Entry point: load → clean → analyze → visualize → train → tune
├── players.json          # Sample (deliberately messy) sports data
├── output/               # Generated charts land here (created on first run)
├── athlete/              # Core domain package
│   ├── athlete.go        # Base Athlete struct + Performer interface
│   └── specialized.go    # Footballer & Cricketer (embed Athlete)
├── preprocess/           # Data cleansing package
│   └── clean.go          # Normalize, validate, sanitize, deduplicate
├── stats/                # Exploratory Data Analysis package
│   └── eda.go            # Describe() summary stats + Frequency() counts
├── viz/                  # Visualization package
│   └── chart.go          # SVG bar chart + residual plot in standalone HTML
└── mlmodel/              # Machine learning package
    └── regression.go     # Linear regression + MSE, residuals, L1 regularization
```

---

## Prerequisites

Choose one path:

- **Docker (with Compose)** — the only requirement for the recommended workflow. [Install Docker](https://docs.docker.com/get-docker/). Compose is bundled with modern Docker Desktop / Engine.
- **Go 1.22+** — if you prefer to run it directly on your machine. [Install Go](https://go.dev/dl/).

---

## Quick Start (Docker Compose)

The simplest way to build and run everything:

```bash
docker compose up --build
```

This builds the image, runs the program (printing EDA and regression results to the console), and copies the generated charts into `./output/` on your machine. Open them in a browser:

- `output/goals_chart.html` — goals per footballer
- `output/residuals.html` — residual diagnostic plot

For a one-off run without recreating the service:

```bash
docker compose run --rm app
```

### How the Compose file works

```yaml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    image: go-sports:latest
    container_name: go-sports
    command: sh -c "./sports-app && cp -f *.html /app/output/ 2>/dev/null || true"
    volumes:
      - ./output:/app/output
    restart: "no"
```

Because go-sports is a **batch job** (it runs, writes its charts, and exits) rather than a long-running server, `restart: "no"` stops Compose from relaunching it. The `volumes` mount maps `./output` on the host into the container, and the `command` copies the generated HTML charts there after the program finishes. The `|| true` guard keeps the run from failing if no chart was produced.

---

## Running with Plain Docker

If you prefer not to use Compose:

```bash
docker build -t go-sports .
docker run --rm -v "$(pwd)/output":/app/output go-sports \
  sh -c "./sports-app && cp -f *.html /app/output/"
```

---

## Running Locally (without Docker)

From the project root:

```bash
go run .
```

This prints results to the console and drops `goals_chart.html` and `residuals.html` right in your project folder. Or build a binary:

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
7. **Fine-tune** — report standard MSE, train an L1-regularized variant for comparison, and write a residual plot to `residuals.html`.

The sample `players.json` is intentionally messy — stray whitespace, mixed casing, a negative goal count, an empty name, and duplicate players — so the cleansing step has real work to do before analysis.

---

## Packages

| Package | Responsibility | Key types / functions |
|---|---|---|
| `athlete` | Domain model | `Athlete`, `Footballer`, `Cricketer`, `Performer` interface |
| `preprocess` | Data cleansing | `CleanFootballers`, `CleanCricketers` |
| `stats` | Exploratory analysis | `Describe` → `Summary`, `Frequency` |
| `viz` | Visualization | `Bar`, `BarChartHTML`, `Point`, `ResidualPlotHTML`, `WriteHTML` |
| `mlmodel` | Machine learning | `Sample`, `LinearModel`, `Normalize`, `Residual` |

### `stats` — Exploratory Data Analysis

`Describe([]float64)` returns a `Summary` with count, min, max, mean, median, and standard deviation — the rough equivalent of pandas' `df.describe()`, built from scratch using the standard `math` and `sort` packages. It copies the input before sorting so the caller's data is never mutated. `Frequency([]string)` returns a `map[string]int` of category counts.

### `viz` — Data Visualization

`BarChartHTML` and `ResidualPlotHTML` assemble inline SVG (using `strings.Builder`) wrapped in standalone HTML pages. Values are scaled proportionally with divide-by-zero guards. `WriteHTML(path, html)` saves to disk. No charting library required.

### `mlmodel` — Linear Regression

A complete simple linear regression (`y = w*x + b`) implemented from scratch, predicting goals from matches played, including fine-tuning tools. See the dedicated sections below.

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
- **Feature scaling matters.** `Normalize` applies min-max scaling so gradient descent converges reliably; unscaled features are the most common cause of a from-scratch model diverging. To predict on a new raw input, scale it the same way before calling `Predict`.
- **Pointer vs. value receivers.** Read-only methods (`Predict`, `Cost`, `RSquared`) use value receivers; the mutating `gradientStep` uses a pointer receiver `(m *LinearModel)` so it can update the parameters in place.

---

## Fine-Tuning: Evaluation & Optimization

The model package also includes tools to evaluate and improve the fit.

| Tool | Implementation | What it does |
|---|---|---|
| **MSE metric** | `MSE(data)` | Standard `1/n` Mean Squared Error for *reporting* (distinct from the `1/2n` `Cost` used for *optimizing*) |
| **Residuals** | `Residuals(data)` → `[]Residual`, `viz.ResidualPlotHTML` | Computes `y − ŷ` per sample and plots it for diagnostics |
| **L1 regularization** | `CostL1`, `TrainL1` | Lasso penalty that discourages large weights and can zero them out |

### Key ideas

- **MSE vs. Cost.** `Cost` uses `1/2n` purely for clean gradients and is what you *optimize*; `MSE` uses the standard `1/n` and is what you *report* and compare across models. Keeping them separate is the conventionally correct setup.
- **Residual plots** are the most useful linear-regression diagnostic. A residual is `actual − predicted`. Random scatter around the zero line means a linear model is appropriate; a curve or funnel shape means the linearity assumption is violated. The plot color-codes over- vs. under-predictions.
- **L1 regularization (Lasso)** adds `lambda × |w|` to the cost; its gradient is `lambda × sign(w)`, a constant push toward zero that can drive weights *exactly* to zero — effectively performing feature selection. The `lambda` hyperparameter controls strength (zero recovers plain regression). The bias is deliberately left unregularized, which is standard practice. With a single well-fit feature the effect is small; it becomes visible with many features, where it zeroes out the irrelevant ones.

---

## Go Concepts Demonstrated

| Concept | Where | What to look for |
|---|---|---|
| Modules & packages | `go.mod`, all folders | Module path `go-sports`; subpackages imported as `go-sports/mlmodel` |
| Structs & struct tags | `athlete/*.go` | `` `json:"name"` `` tags map struct fields to JSON keys |
| Embedding (composition) | `specialized.go` | `Footballer` embeds `Athlete` and inherits its fields/methods |
| Field & method promotion | `main.go` | `footballer.Name` and `footballer.Describe()` come from `Athlete` |
| Interfaces (polymorphism) | `athlete.go` | `Performer` interface; satisfied implicitly by any type with `Stats()` |
| Methods & receivers | `athlete.go`, `eda.go`, `regression.go` | Value receivers for reads; pointer receivers for mutating steps |
| JSON marshal/unmarshal | `main.go` | `encoding/json` for decoding |
| Error handling | `main.go` | The idiomatic `if err != nil` pattern |
| Maps as sets / counters | `clean.go`, `eda.go` | `map[string]bool` to dedupe; `map[string]int` to count |
| Regular expressions | `preprocess/clean.go` | `regexp` collapses repeated whitespace |
| Math & sorting | `stats/eda.go`, `regression.go` | `math.Sqrt`, `math.Abs`, `sort.Float64s` |
| Efficient string building | `viz/chart.go` | `strings.Builder` to assemble SVG |
| Multiple return values | `mlmodel/regression.go` | `Normalize` returns scaled data plus min and span |
| Closures | `viz/chart.go` | `sx`/`sy` coordinate-mapping functions capture plot dimensions |
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

===== Fine-Tuning =====
MSE (no regularization): 0.0024
L1 (lambda=0.10): weight=0.452 bias=0.041  MSE=0.0031
Residual plot written to residuals.html
```

Exact numbers vary, but cost should fall sharply, R² should be close to 1.0 on this clean sample, and the L1 model's weight should be slightly smaller than the unregularized one.

---

## Next Steps

Ideas to extend the project as you keep learning:

- Add a **cost-history chart** to `viz` so you can see the model converge.
- Extend to **multiple features** (multivariate regression) to make L1's feature-selection effect visible.
- Add **L2 regularization (Ridge)** and compare it with L1.
- Add a **train/test split** to evaluate on unseen data.
- Make cleansing **configurable** with a `CleanOptions` struct.
- Read from **CSV** instead of JSON to practice another format.
- Expose the cleaned data, charts, and predictions over an **HTTP API** (then add a `ports:` mapping and `restart: unless-stopped` to the Compose file).
- Add **unit tests** (`go test`) for the cleansing, stats, and model functions.

---

## License

This is a learning project — use it freely.
