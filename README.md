# go-sports

A small, beginner-friendly Go project for learning the full data workflow — **loading, cleansing, exploratory analysis, visualization, and machine learning from scratch** — built around a sports theme and packaged to run in Docker. It demonstrates Go modules, struct embedding (Go's take on inheritance), interfaces, a preprocessing pipeline, descriptive statistics, dependency-free SVG charting, and three ML model families implemented without any ML libraries: linear regression, logistic regression, and decision trees / random forests — with feature importance and model explainability.

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
  - [Feature Importance & Explainability](#feature-importance--explainability)
- [Go Concepts Demonstrated](#go-concepts-demonstrated)
- [Sample Output](#sample-output)
- [Interpreting the Output](#interpreting-the-output)
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
- **Feature importance & explainability**: forest-wide Gini importance plus per-prediction decision-path explanations.
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
    ├── tree.go           # Decision tree: Gini, splitting, build, importance, Explain
    └── forest.go         # Random forest: bagging, subsampling, voting, importance
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
- `output/importance.html` — feature importance bar chart (if enabled)

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

This prints results to the console and drops the chart HTML files right in your project folder. Or build a binary:

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
5. **Visualize** — the `viz` package builds an SVG bar chart.
6. **Linear regression** — fit, evaluate, predict; then fine-tune with MSE, a residual plot, and L1 regularization.
7. **Logistic regression** — classify a binary label with the sigmoid and log-loss; report metrics, decision boundary, and a softmax example.
8. **Trees & forest** — train a decision tree and a random forest, then predict and score.
9. **Importance & explainability** — rank features by Gini importance and trace a single prediction's decision path.

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
| `tree` | Tree-based models | `Row`, `DecisionTree`, `RandomForest`, `Step` |

---

## Machine Learning Models

All models are implemented from scratch with no external ML dependencies — the point is to see the math, not hide it behind a library.

### Linear Regression

Simple linear regression (`y = w*x + b`) predicting goals from matches played, covering all six standard stages: model & data (`Sample`, `LinearModel`), hypothesis (`Predict`), cost (`Cost`, MSE with a `1/2n` factor for clean gradients), gradient descent (`gradientStep`), training (`Train`, returns cost history), and evaluation (`RSquared`, `RMSE`). Feature scaling via `Normalize` keeps gradient descent stable; read-only methods use value receivers while the mutating `gradientStep` uses a pointer receiver.

### Fine-Tuning: Evaluation & Optimization

A reporting-grade `MSE` (`1/n`, distinct from the `1/2n` optimization `Cost`), `Residuals` plotted via `viz.ResidualPlotHTML` for diagnostics, and L1/Lasso regularization (`CostL1`, `TrainL1`) that adds `lambda × |w|` to the cost — its gradient `lambda × sign(w)` is a constant push toward zero that can perform feature selection. The bias is left unregularized.

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

Log loss is used because MSE is non-convex with the sigmoid; precision/recall/F1 matter because accuracy misleads on imbalanced classes; the decision boundary for threshold 0.5 reduces to `x = −b/w`; softmax outputs a probability distribution summing to 1.

### Decision Trees & Random Forests

Tree-based classification predicting the top-scorer label from `[matches, goals]`.

| Concept | Implementation | What it does |
|---|---|---|
| **Data** | `Row` | Feature vector plus a binary label |
| **Impurity** | `gini(rows)` | 0 = pure, ~0.5 = maximally mixed for two classes |
| **Splitting** | `bestSplit`, `split` | Finds the feature/threshold with the largest information gain |
| **Tree** | `DecisionTree`, `Node` | Recursive `*Node` structure; leaves predict by majority |
| **Ensemble** | `RandomForest` | Many trees on bootstrap samples |
| **Bagging** | `bootstrap(rows)` | Resamples rows with replacement per tree |
| **Feature subsampling** | `featuresPerSplit` | Each split considers `sqrt(features)` random features |
| **Voting** | `Predict` | Majority vote across all trees |

A single deep tree overfits; a forest reduces this with two sources of randomness — bagging and feature subsampling — which decorrelate the trees so averaging their votes generalizes far better.

### Feature Importance & Explainability

These answer two different questions: importance is **global** (which features matter across the whole model), while an explanation is **local** (why this one prediction came out as it did).

| Tool | Implementation | What it does |
|---|---|---|
| **Tree importance** | `DecisionTree.FeatureImportance` | Sums `samples × gain` for each feature over all splits |
| **Forest importance** | `RandomForest.FeatureImportance` | Averages per-tree importance, normalized to sum to 1 |
| **Decision path** | `DecisionTree.Explain` → `[]Step` | Traces root-to-leaf, recording each threshold test |

**Gini importance** credits a feature with the impurity decrease of every split that uses it, weighted by how many samples pass through that node — so features that produce big, clean separations near the root score highest. The forest averages this across all trees and normalizes, so each value reads as that feature's share of the decision-making.

**The decision path** is the most faithful explanation a tree can give because it *is* the computation — no post-hoc approximation. `Explain` records each threshold test from root to leaf, producing a literal chain like "matches (33.0) > 25.0 → class 1." This is why single trees are prized for interpretability: the model and its explanation are the same object.

> **Caveat:** Gini importance is biased toward high-cardinality features (more distinct values give more candidate thresholds). **Permutation importance** — shuffling one feature and measuring the accuracy drop — avoids that bias and is a natural next addition.

To record importance data, each `Node` stores `NSamples`, `Impurity`, and `Gain` during training.

---

## Go Concepts Demonstrated

| Concept | Where | What to look for |
|---|---|---|
| Modules & packages | `go.mod`, all folders | Module path `go-sports`; subpackages like `go-sports/tree` |
| Structs & struct tags | `athlete/*.go` | `` `json:"name"` `` tags map struct fields to JSON keys |
| Embedding (composition) | `specialized.go` | `Footballer` embeds `Athlete` |
| Interfaces (polymorphism) | `athlete.go` | `Performer` satisfied implicitly |
| Methods & receivers | `regression.go`, `logistic.go` | Value receivers for reads; pointer receivers for mutation |
| Recursion & pointer structures | `tree/tree.go` | `*Node` tree built, traversed, and explained recursively |
| Closures | `tree/tree.go`, `viz/chart.go` | Recursive `walk` closure for importance; `sx`/`sy` mappers |
| Encapsulation | `tree/tree.go` | Unexported `featuresPerSplit` hides forest internals |
| Randomness | `tree/forest.go` | `math/rand` for bootstrap sampling and feature shuffling |
| Maps as counters | `clean.go`, `eda.go`, `tree.go` | Label counts, dedup sets, vote tallies |
| JSON marshal/unmarshal | `main.go` | `encoding/json` for decoding |
| Error handling | `main.go` | The idiomatic `if err != nil` pattern |
| Regular expressions | `preprocess/clean.go` | `regexp` collapses repeated whitespace |
| Math & sorting | `stats`, `mlmodel`, `tree` | `math.Exp`, `math.Sqrt`, `sort.Float64s` |
| Efficient string building | `viz/chart.go` | `strings.Builder` to assemble SVG |
| Multiple return values | `regression.go`, `tree.go` | `Normalize`, `Explain` |

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
===== Linear Regression: predict goals from matches =====
Learned model: goals = 0.470 * matches + 0.030
  R²   : 0.998
  RMSE : 0.350

===== Logistic Regression: is this player a top scorer? =====
  Accuracy : 1.000  Precision: 1.000  Recall: 1.000  F1: 1.000
  Decision boundary at x = 0.421

===== Decision Tree & Random Forest =====
Decision tree predicts class 1 for [33 16]
Random forest predicts class 1 for [33 16]
Forest training accuracy: 1.000

===== Feature Importance & Explainability =====
Feature importance (forest, normalized):
   matches : 0.521
   goals   : 0.479

Why the tree predicted class 1 for [33 16]:
   matches (33.0) > 25.0
```

---

## Interpreting the Output

### Feature importance

```
matches : 0.521
goals   : 0.479
```

These are the forest's normalized Gini-importance scores, so they **sum to 1.0** and each reads as that feature's share of the total decision-making. Here `matches` (52%) and `goals` (48%) are almost equally important. That near-even split is expected when the two features are **highly correlated** — players with more matches tend to have more goals — so across 50 trees with random feature subsets, the credit gets shared roughly evenly rather than concentrating in one feature. Neither dominates.

### The decision path

```
Why the tree predicted class 1 for [33 16]:
   matches (33.0) > 25.0
```

This is the **literal path** the single decision tree walked for the player `[33 matches, 16 goals]`. It needed only **one test**: is `matches` (33.0) greater than 25.0? Yes — so it followed that branch to a leaf labeled class 1 (top scorer). The path is a single line because the data separates cleanly at "matches > 25," so the root split alone is decisive and the tree never had to examine `goals`.

### Why they seem to disagree (but don't)

The forest says both features matter (0.52 / 0.48), yet this one tree's explanation used only `matches`. That's not a contradiction. Global importance is **averaged over 50 bootstrapped trees**, each trained on different resampled data and feature subsets — many of those trees split on `goals`, which is how `goals` earns its 0.48 share. The explanation, by contrast, is **one specific tree's** path. Global (importance) and local (path) views answer different questions and are expected to look different.

---

## Next Steps

Ideas to extend the project as you keep learning:

- Add **permutation importance** as a less-biased alternative to Gini importance.
- Visualize the **full tree structure** (not just one path) as an SVG.
- Add a **cost-history chart** to `viz` so you can watch models converge.
- Extend regression to **multiple features** to make L1's feature-selection effect visible.
- Add **L2 regularization (Ridge)** and compare it with L1.
- Add a **train/test split**, **cross-validation**, and **out-of-bag (OOB) error**.
- Read from **CSV** instead of JSON to practice another format.
- Expose the cleaned data, charts, and predictions over an **HTTP API**.
- Add **unit tests** (`go test`) for the cleansing, stats, and model functions.

---

## License

This is a learning project — use it freely.
