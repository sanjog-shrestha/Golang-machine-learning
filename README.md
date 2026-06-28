# go-sports

A small, beginner-friendly Go project for learning the full data workflow — **loading, cleansing, exploratory analysis, visualization, and machine learning from scratch** — built around a sports theme and packaged to run in Docker. It demonstrates Go modules, struct embedding (Go's take on inheritance), interfaces, a preprocessing pipeline, descriptive statistics, dependency-free SVG charting, and three ML model families implemented without any ML libraries: linear regression, logistic regression, and decision trees / random forests.

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
- [Machine Learning Models](#machine-learning-models)
  - [Linear Regression](#linear-regression)
  - [Fine-Tuning: Evaluation & Optimization](#fine-tuning-evaluation--optimization)
  - [Logistic Regression](#logistic-regression)
  - [Decision Trees & Random Forests](#decision-trees--random-forests)
- [Go Concepts Demonstrated](#go-concepts-demonstrated)
- [Sample Output](#sample-output)
- [Next Steps](#next-steps)

---

## Features

- Load sports data from a JSON file into typed Go structs.
- Organize code into **modules and packages** (`athlete`, `preprocess`, `stats`, `viz`, `mlmodel`, `tree`).
- Model data with **struct embedding** (composition instead of inheritance).
- Use **interfaces** for polymorphism across different athlete types.
- A **data cleansing & preprocessing** pipeline that normalizes, validates, sanitizes, and deduplicates messy input.
- **Exploratory Data Analysis (EDA)**: descriptive statistics plus categorical frequency counts.
- **Data visualization**: self-contained SVG charts (bar chart + residual plot) written to HTML files — no external dependencies.
- **Linear regression from scratch**: hypothesis, cost, gradient descent, training, prediction, evaluation, MSE, residuals, and L1 regularization.
- **Logistic regression from scratch**: sigmoid, log-loss, gradient descent, classification metrics, decision boundaries, and softmax for multi-class.
- **Decision trees & random forests from scratch**: Gini impurity, recursive splitting, bagging, feature subsampling, and majority voting.
- A **multi-stage Docker build** plus a **Docker Compose** workflow that writes charts back to your host.

---

## Project Structure

```
go-sports/
├── Dockerfile            # Multi-stage build (compile in Go image, run on Alpine)
├── docker-compose.yml    # One-command build/run; mounts ./output for charts
├── go.mod                # Module definition: module path "go-sports"
├── main.go               # Entry point: load → clean → analyze → visualize → ML
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
├── mlmodel/              # Regression & classification package
│   ├── regression.go     # Linear regression + MSE, residuals, L1 regularization
│   └── logistic.go       # Logistic regression + sigmoid, metrics, softmax
└── tree/                 # Tree-based models package
    ├── tree.go           # Decision tree: Gini, splitting, recursive build
    └── forest.go         # Random forest: bagging, feature subsampling, voting
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

This builds the image, runs the program (printing EDA and all model results to the console), and copies the generated charts into `./output/` on your machine. Open them in a browser:

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
6. **Linear regression** — fit with gradient descent, evaluate, predict; then fine-tune with MSE, a residual plot, and L1 regularization.
7. **Logistic regression** — classify a binary label with the sigmoid and log-loss; report metrics, decision boundary, and a softmax example.
8. **Trees & forest** — train a decision tree and a random forest, then predict and score.

The sample `players.json` is intentionally messy so the cleansing step has real work to do before analysis.

---

## Packages

| Package | Responsibility | Key types / functions |
|---|---|---|
| `athlete` | Domain model | `Athlete`, `Footballer`, `Cricketer`, `Performer` interface |
| `preprocess` | Data cleansing | `CleanFootballers`, `CleanCricketers` |
| `stats` | Exploratory analysis | `Describe` → `Summary`, `Frequency` |
| `viz` | Visualization | `Bar`, `BarChartHTML`, `Point`, `ResidualPlotHTML`, `WriteHTML` |
| `mlmodel` | Regression & classification | `LinearModel`, `LogisticModel`, `Normalize`, `Sigmoid`, `Softmax`, `Metrics` |
| `tree` | Tree-based models | `Row`, `DecisionTree`, `RandomForest` |

---

## Machine Learning Models

All models are implemented from scratch with no external ML dependencies — the point is to see the math, not hide it behind a library.

### Linear Regression

Simple linear regression (`y = w*x + b`) predicting goals from matches played, covering all six standard stages.

| Stage | Implementation | What it does |
|---|---|---|
| **1. Model & data** | `Sample`, `LinearModel` | A `Sample` is one (feature, target) pair; `LinearModel` holds `Weight` and `Bias` |
| **2. Hypothesis** | `Predict(x)` | Computes `h(x) = w*x + b` |
| **3. Cost** | `Cost(data)` | Mean Squared Error: `J = (1/2n) Σ (h(xᵢ) − yᵢ)²` |
| **4. Gradient descent** | `gradientStep(data, lr)` | Partial derivatives nudge parameters downhill |
| **5. Training** | `Train(data, lr, epochs)` | Repeats steps, returns cost history |
| **6. Prediction & evaluation** | `Predict`, `RSquared`, `RMSE` | Forecasts and reports goodness of fit |

Key points: MSE squares errors (penalizing big ones, staying positive); the learning rate sets step size; **feature scaling** via `Normalize` keeps gradient descent stable; read-only methods use value receivers while the mutating `gradientStep` uses a pointer receiver.

### Fine-Tuning: Evaluation & Optimization

| Tool | Implementation | What it does |
|---|---|---|
| **MSE metric** | `MSE(data)` | Standard `1/n` MSE for *reporting* (distinct from the `1/2n` `Cost` used for *optimizing*) |
| **Residuals** | `Residuals(data)`, `viz.ResidualPlotHTML` | Computes `y − ŷ` and plots it for diagnostics |
| **L1 regularization** | `CostL1`, `TrainL1` | Lasso penalty that discourages large weights and can zero them out |

A residual plot showing random scatter around zero confirms a linear model fits; a curve or funnel means it doesn't. L1 adds `lambda × |w|` to the cost; its gradient `lambda × sign(w)` is a constant push toward zero that can perform feature selection. The bias is left unregularized.

### Logistic Regression

Binary classification — "is this player a top scorer?" — covering all seven stages.

| Stage | Implementation | What it does |
|---|---|---|
| **1. Model & data** | `LabeledSample`, `LogisticModel` | Same `Weight`/`Bias` shape; target is a 0/1 label |
| **2. Sigmoid** | `Sigmoid(z)`, `Probability(x)` | Squashes `w*x + b` into a `(0,1)` probability |
| **3. Cost** | `Cost(data)` | Binary cross-entropy (log loss), not MSE |
| **4. Gradient descent** | `gradientStep(data, lr)` | Gradient simplifies to the same `(ŷ − y)·x` form |
| **5. Prediction & evaluation** | `Classify`, `Evaluate` → `Metrics` | Confusion matrix, accuracy, precision, recall, F1 |
| **6. Decision boundary** | `DecisionBoundary(threshold)` | The feature value where `P(y=1) = threshold` |
| **7. Beyond binary** | `Softmax(scores)` | Multi-class generalization of the sigmoid |

Key points: log loss is used because MSE is non-convex with the sigmoid; precision/recall/F1 matter because accuracy misleads on imbalanced classes; the decision boundary for threshold 0.5 reduces to `x = −b/w`; softmax outputs a probability distribution summing to 1, with a max-subtraction trick for numerical stability.

### Decision Trees & Random Forests

Tree-based classification predicting the top-scorer label from `[matches, goals]`.

| Concept | Implementation | What it does |
|---|---|---|
| **Data** | `Row` | Feature vector plus a binary label |
| **Impurity** | `gini(rows)` | 0 = pure, ~0.5 = maximally mixed for two classes |
| **Splitting** | `bestSplit`, `split` | Finds the feature/threshold with the largest information gain |
| **Tree** | `DecisionTree`, `Node` | Recursive `*Node` structure; leaves predict by majority |
| **Stopping** | `MaxDepth`, `MinLeafSize` | Halts growth to limit overfitting |
| **Ensemble** | `RandomForest` | Many trees trained on bootstrap samples |
| **Bagging** | `bootstrap(rows)` | Resamples rows with replacement per tree |
| **Feature subsampling** | `featuresPerSplit` | Each split considers `sqrt(features)` random features |
| **Voting** | `Predict` | Majority vote across all trees |

Key points: a single deep tree overfits by memorizing the data. A random forest reduces this with two sources of randomness — bagging (each tree sees a different resample) and feature subsampling (trees split on different features) — which decorrelates the trees so averaging their votes generalizes far better. The tree is built from `*Node` pointers because a recursive, dynamically-shaped structure can't be a plain value; `featuresPerSplit` is unexported, so a lone tree defaults to all features while the forest sets it internally.

---

## Go Concepts Demonstrated

| Concept | Where | What to look for |
|---|---|---|
| Modules & packages | `go.mod`, all folders | Module path `go-sports`; subpackages like `go-sports/tree` |
| Structs & struct tags | `athlete/*.go` | `` `json:"name"` `` tags map struct fields to JSON keys |
| Embedding (composition) | `specialized.go` | `Footballer` embeds `Athlete` |
| Interfaces (polymorphism) | `athlete.go` | `Performer` satisfied implicitly |
| Methods & receivers | `regression.go`, `logistic.go` | Value receivers for reads; pointer receivers for mutation |
| Recursion & pointer structures | `tree/tree.go` | `*Node` tree built and traversed recursively |
| Encapsulation | `tree/tree.go` | Unexported `featuresPerSplit` hides forest internals |
| Randomness | `tree/forest.go` | `math/rand` for bootstrap sampling and feature shuffling |
| Maps as counters | `clean.go`, `eda.go`, `tree.go` | Label counts, dedup sets, vote tallies |
| JSON marshal/unmarshal | `main.go` | `encoding/json` for decoding |
| Error handling | `main.go` | The idiomatic `if err != nil` pattern |
| Regular expressions | `preprocess/clean.go` | `regexp` collapses repeated whitespace |
| Math & sorting | `stats`, `mlmodel`, `tree` | `math.Exp`, `math.Sqrt`, `sort.Float64s` |
| Efficient string building | `viz/chart.go` | `strings.Builder` to assemble SVG |
| Multiple return values | `regression.go`, `logistic.go` | `Normalize`, `DecisionBoundary` |
| Closures | `viz/chart.go` | `sx`/`sy` coordinate mappers capture plot dimensions |

### A note on "inheritance"

Go has **no classical inheritance**. Instead it uses **composition** via struct embedding (a `Footballer` *has an* `Athlete`) and **interfaces** for shared behavior. This is more explicit than subclassing and avoids fragile deep class hierarchies.

### Data cleansing stages

- **Normalization** — trim/collapse whitespace and unify casing (`MESSI` → `Messi`).
- **Validation** — drop records missing required fields (e.g. empty name).
- **Sanitization** — fix out-of-range values (a negative goal count is clamped to `0`).
- **Deduplication** — collapse repeats using a composite `name|team` key.

---

## Sample Output

Console output follows this shape (exact numbers vary):

```
===== EDA: Football Goals =====
--- Goals ---
  Count : 2
  Mean  : 7.50
  ...

===== Linear Regression: predict goals from matches =====
Learned model: goals = 0.470 * matches + 0.030
  R²   : 0.998
  RMSE : 0.350

===== Fine-Tuning =====
MSE (no regularization): 0.0024
L1 (lambda=0.10): weight=0.452 bias=0.041  MSE=0.0031
Residual plot written to residuals.html

===== Logistic Regression: is this player a top scorer? =====
Log-loss: start=0.6931 end=0.0456
  Confusion: TP=3 TN=3 FP=0 FN=0
  Accuracy : 1.000
  Precision: 1.000
  Recall   : 1.000
  F1 Score : 1.000
  Decision boundary at x = 0.421
Player at x=0.60: P(top scorer)=0.912 -> class 1
Multi-class softmax [2 1 0.1] -> 0.659, 0.242, 0.099

===== Decision Tree & Random Forest =====
Decision tree predicts class 1 for [33 16]
Random forest predicts class 1 for [33 16]
Forest training accuracy: 1.000
```

---

## Next Steps

Ideas to extend the project as you keep learning:

- Add **feature importance** to the random forest (which features drive the splits).
- Add a **cost-history chart** to `viz` so you can watch models converge.
- Extend regression to **multiple features** to make L1's feature-selection effect visible.
- Add **L2 regularization (Ridge)** and compare it with L1.
- Add a **train/test split** and **cross-validation** to evaluate on unseen data.
- Add **out-of-bag (OOB) error** estimation to the forest.
- Read from **CSV** instead of JSON to practice another format.
- Expose the cleaned data, charts, and predictions over an **HTTP API**.
- Add **unit tests** (`go test`) for the cleansing, stats, and model functions.

---

## License

This is a learning project — use it freely.
