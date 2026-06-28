# go-sports

A small, beginner-friendly Go project for learning the full data workflow — **loading, cleansing, exploratory analysis, visualization, and from-scratch machine learning models (linear *and* logistic regression)** — built around a sports theme and packaged to run in Docker. It demonstrates Go modules, struct embedding (Go's take on inheritance), interfaces, a preprocessing pipeline, descriptive statistics, dependency-free SVG charting, and regression/classification models built without any ML libraries.

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
- [Logistic Regression for Binary Classification](#logistic-regression-for-binary-classification)
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
- **Exploratory Data Analysis (EDA)**: descriptive statistics plus categorical frequency counts.
- **Data visualization**: self-contained SVG charts (bar chart + residual plot) written to HTML files — no external dependencies.
- **Linear regression from scratch**: hypothesis, cost, gradient descent, training, prediction, and evaluation.
- **Model fine-tuning**: standard MSE metric, residual diagnostics, and L1 (Lasso) regularization.
- **Logistic regression from scratch**: sigmoid, log-loss, gradient descent, classification metrics, decision boundaries, and a softmax for multi-class.
- A **multi-stage Docker build** plus a **Docker Compose** workflow that writes charts back to your host.

---

## Project Structure

```
go-sports/
├── Dockerfile            # Multi-stage build (compile in Go image, run on Alpine)
├── docker-compose.yml    # One-command build/run; mounts ./output for charts
├── go.mod                # Module definition: module path "go-sports"
├── main.go               # Entry point: load → clean → analyze → visualize → train → tune → classify
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
    ├── regression.go     # Linear regression + MSE, residuals, L1 regularization
    └── logistic.go       # Logistic regression: sigmoid, log-loss, metrics, softmax
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

This builds the image, runs the program (printing EDA, regression, and classification results to the console), and copies the generated charts into `./output/` on your machine. Open them in a browser:

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
4. **Analyze (EDA)** — the `stats` package computes descriptive statistics and frequency counts.
5. **Visualize** — the `viz` package builds an SVG bar chart and writes it to `goals_chart.html`.
6. **Train & predict (linear)** — the `mlmodel` package fits a linear regression model, evaluates it, and makes a prediction.
7. **Fine-tune** — report standard MSE, train an L1-regularized variant, and write a residual plot.
8. **Classify (logistic)** — fit a logistic regression model, evaluate it with a confusion matrix, find its decision boundary, and demonstrate softmax for multi-class.

The sample `players.json` is intentionally messy — stray whitespace, mixed casing, a negative goal count, an empty name, and duplicate players — so the cleansing step has real work to do before analysis.

---

## Packages

| Package | Responsibility | Key types / functions |
|---|---|---|
| `athlete` | Domain model | `Athlete`, `Footballer`, `Cricketer`, `Performer` interface |
| `preprocess` | Data cleansing | `CleanFootballers`, `CleanCricketers` |
| `stats` | Exploratory analysis | `Describe` → `Summary`, `Frequency` |
| `viz` | Visualization | `Bar`, `BarChartHTML`, `Point`, `ResidualPlotHTML`, `WriteHTML` |
| `mlmodel` | Machine learning | `LinearModel`, `LogisticModel`, `Normalize`, `Sigmoid`, `Softmax`, `Metrics` |

---

## The Linear Regression Model

The `mlmodel` package implements all six stages of building a regression model, with no external ML dependencies.

| Stage | Implementation | What it does |
|---|---|---|
| **1. Model & data** | `Sample`, `LinearModel` | A `Sample` is one (feature, target) pair; `LinearModel` holds `Weight` (slope) and `Bias` (intercept) |
| **2. Hypothesis** | `Predict(x)` | Computes `h(x) = w*x + b` |
| **3. Cost function** | `Cost(data)` | Mean Squared Error: `J = (1/2n) Σ (h(xᵢ) − yᵢ)²` |
| **4. Gradient descent** | `gradientStep(data, lr)` | Computes partial derivatives and nudges parameters downhill |
| **5. Training** | `Train(data, lr, epochs)` | Repeats gradient steps, returns the cost history |
| **6. Prediction & evaluation** | `Predict`, `RSquared`, `RMSE`, `Summary` | Forecasts new values and reports goodness of fit |

Key ideas: **MSE** penalizes large errors; **gradient descent** moves parameters opposite the gradient with a **learning rate** over many **epochs**; **feature scaling** (`Normalize`) keeps it from diverging; and read-only methods use value receivers while the mutating `gradientStep` uses a pointer receiver.

---

## Fine-Tuning: Evaluation & Optimization

| Tool | Implementation | What it does |
|---|---|---|
| **MSE metric** | `MSE(data)` | Standard `1/n` MSE for *reporting* (distinct from the `1/2n` `Cost` used for *optimizing*) |
| **Residuals** | `Residuals(data)`, `viz.ResidualPlotHTML` | Computes `y − ŷ` per sample and plots it for diagnostics |
| **L1 regularization** | `CostL1`, `TrainL1` | Lasso penalty that discourages large weights and can zero them out |

Key ideas: **MSE vs. Cost** are deliberately separate (report vs. optimize); **residual plots** reveal whether a linear model is appropriate (random scatter = good, patterns = bad); **L1 (Lasso)** adds `lambda × |w|` whose gradient `lambda × sign(w)` can drive weights exactly to zero for feature selection, with the bias left unregularized.

---

## Logistic Regression for Binary Classification

The `mlmodel` package also implements logistic regression from scratch, predicting whether a player is a **top scorer** (1) or not (0). It covers all seven stages.

| Stage | Implementation | What it does |
|---|---|---|
| **1. Model & data** | `LabeledSample`, `LogisticModel` | Same `Weight`/`Bias` shape as linear, but output is a probability and labels are 0/1 |
| **2. Sigmoid** | `Sigmoid(z)`, `Probability(x)` | Squashes `w*x + b` into `(0,1)` so it reads as `P(y=1)` |
| **3. Cost function** | `Cost(data)` | Binary cross-entropy (log loss), which is convex — MSE is not used here |
| **4. Gradient descent** | `gradientStep`, `Train` | Same gradient form as linear regression: `(ŷ − y)·x` |
| **5. Prediction & evaluation** | `Classify`, `Evaluate` → `Metrics` | Confusion matrix with accuracy, precision, recall, F1 |
| **6. Decision boundary** | `DecisionBoundary(threshold)` | Feature value where `P(y=1) = threshold`; for 0.5 it's `−b/w` |
| **7. Beyond binary** | `Softmax(scores)` | Multi-class generalization: a probability distribution summing to 1 |

### Key ideas

- **Sigmoid** turns a linear score into a probability — the defining step that makes classification possible.
- **Log loss, not MSE.** MSE with a sigmoid is non-convex (local minima); cross-entropy is convex, so gradient descent finds the global minimum. A small epsilon guards against `log(0)`.
- **Identical gradient.** After differentiating log loss through the sigmoid, the update simplifies to the same `(prediction − actual) × feature` form as linear regression — only the prediction and cost differ.
- **Evaluation beyond accuracy.** A confusion matrix yields **precision** (correctness of positive predictions), **recall** (coverage of actual positives), and **F1** (their harmonic mean) — essential when classes are imbalanced.
- **Decision boundary** is where the model is exactly undecided. In 1-D it's a point (`−b/w`); with two features a line; in higher dimensions a hyperplane.
- **Softmax** generalizes the sigmoid to many classes, outputting a probability distribution. Subtracting the max score before exponentiating is a numerical-stability trick that avoids overflow.

---

## Go Concepts Demonstrated

| Concept | Where | What to look for |
|---|---|---|
| Modules & packages | `go.mod`, all folders | Module path `go-sports`; subpackages imported as `go-sports/mlmodel` |
| Structs & struct tags | `athlete/*.go` | `` `json:"name"` `` tags map struct fields to JSON keys |
| Embedding (composition) | `specialized.go` | `Footballer` embeds `Athlete` and inherits its fields/methods |
| Field & method promotion | `main.go` | `footballer.Name` and `footballer.Describe()` come from `Athlete` |
| Interfaces (polymorphism) | `athlete.go` | `Performer` interface; satisfied implicitly by any type with `Stats()` |
| Methods & receivers | `regression.go`, `logistic.go` | Value receivers for reads; pointer receivers for mutating steps |
| JSON marshal/unmarshal | `main.go` | `encoding/json` for decoding |
| Error handling | `main.go` | The idiomatic `if err != nil` pattern |
| Maps as sets / counters | `clean.go`, `eda.go` | `map[string]bool` to dedupe; `map[string]int` to count |
| Regular expressions | `preprocess/clean.go` | `regexp` collapses repeated whitespace |
| Math functions | `regression.go`, `logistic.go` | `math.Sqrt`, `math.Abs`, `math.Exp`, `math.Log` |
| Switch statements | `logistic.go` | Confusion-matrix tallying and the `sign` helper |
| Efficient string building | `viz/chart.go` | `strings.Builder` to assemble SVG |
| Multiple return values | `mlmodel/*.go` | `Normalize` returns scaled data + min + span; `DecisionBoundary` returns value + ok |
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

===== Logistic Regression: is this player a top scorer? =====
Log-loss: start=0.6931 end=0.0488
  Confusion: TP=3 TN=3 FP=0 FN=0
  Accuracy : 1.000
  Precision: 1.000
  Recall   : 1.000
  F1 Score : 1.000
  Decision boundary at x = 0.421

Player at x=0.60: P(top scorer)=0.973 -> class 1

Multi-class softmax [2 1 0.1] -> 0.659, 0.242, 0.099
```

Exact numbers vary, but loss should fall sharply, the linear model's R² should be near 1.0, and the classifier should separate the two classes cleanly on this simple sample.

---

## Next Steps

Ideas to extend the project as you keep learning:

- Add a **sigmoid-curve chart** with the decision boundary marked to `viz`.
- Add a **cost-history chart** so you can see models converge.
- Extend to **multiple features** (multivariate regression / classification).
- Implement full **multinomial logistic regression** using the softmax.
- Add a **train/test split** to evaluate on unseen data.
- Add **L2 regularization (Ridge)** and compare it with L1.
- Read from **CSV** instead of JSON to practice another format.
- Expose the models over an **HTTP API** (then add a `ports:` mapping and `restart: unless-stopped` to the Compose file).
- Add **unit tests** (`go test`) for the cleansing, stats, and model functions.

---

## License

This is a learning project — use it freely.
