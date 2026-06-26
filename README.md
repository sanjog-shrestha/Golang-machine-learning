# go-sports

A small, beginner-friendly Go project for learning **data loading, handling, and cleansing** with a sports theme — packaged to run in Docker. It demonstrates Go modules, struct embedding (Go's take on inheritance), interfaces, and a data preprocessing pipeline.

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Quick Start (Docker)](#quick-start-docker)
- [Running Locally (without Docker)](#running-locally-without-docker)
- [How It Works](#how-it-works)
- [Go Concepts Demonstrated](#go-concepts-demonstrated)
- [Sample Output](#sample-output)
- [Next Steps](#next-steps)

---

## Features

- Load sports data from a JSON file into typed Go structs.
- Organize code into **modules and packages** (`athlete`, `preprocess`).
- Model data with **struct embedding** (composition instead of inheritance).
- Use **interfaces** for polymorphism across different athlete types.
- A **data cleansing & preprocessing** pipeline that normalizes, validates, sanitizes, and deduplicates messy input.
- A **multi-stage Docker build** producing a small Alpine-based runtime image.

---

## Project Structure

```
go-sports/
├── Dockerfile            # Multi-stage build (compile in Go image, run on Alpine)
├── go.mod                # Module definition: module path "go-sports"
├── main.go               # Entry point: load → clean → process → print
├── players.json          # Sample (deliberately messy) sports data
├── athlete/              # Core domain package
│   ├── athlete.go        # Base Athlete struct + Performer interface
│   └── specialized.go    # Footballer & Cricketer (embed Athlete)
└── preprocess/           # Data cleansing package
    └── clean.go          # Normalize, validate, sanitize, deduplicate
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

Run the container (the `--rm` flag removes it after it exits):

```bash
docker run --rm go-sports
```

That's it — the program loads the data, cleans it, and prints the results.

---

## Running Locally (without Docker)

From the project root:

```bash
go run .
```

Or build a binary and run it:

```bash
go build -o sports-app .
./sports-app
```

---

## How It Works

The program runs a simple pipeline in `main.go`:

1. **Load** — `os.ReadFile` reads `players.json` as raw bytes.
2. **Parse** — `json.Unmarshal` decodes the JSON into a `Roster` struct.
3. **Cleanse** — the `preprocess` package normalizes text, drops invalid rows, sanitizes numbers, and removes duplicates.
4. **Process** — cleaned records are gathered into a slice of the `Performer` interface and their behaviors are called.
5. **Report** — counts and descriptions are printed before and after cleaning, so you can see what the pipeline changed.

The sample `players.json` is intentionally messy — stray whitespace, mixed casing, a negative goal count, an empty name, and duplicate players — so the cleansing step has real work to do.

---

## Go Concepts Demonstrated

| Concept | Where | What to look for |
|---|---|---|
| Modules & packages | `go.mod`, all folders | Module path `go-sports`; subpackages imported as `go-sports/athlete` |
| Structs & struct tags | `athlete/*.go` | `` `json:"name"` `` tags map struct fields to JSON keys |
| Embedding (composition) | `specialized.go` | `Footballer` embeds `Athlete` and inherits its fields/methods |
| Field & method promotion | `main.go` | `footballer.Name` and `footballer.Describe()` come from `Athlete` |
| Interfaces (polymorphism) | `athlete.go` | `Performer` interface; satisfied implicitly by any type with `Stats()` |
| Methods & receivers | `athlete.go` | `func (a Athlete) Describe()` — `a` is the receiver |
| JSON marshal/unmarshal | `main.go` | `encoding/json` for decoding and encoding |
| Error handling | `main.go` | The idiomatic `if err != nil` pattern |
| Maps as sets | `preprocess/clean.go` | `map[string]bool` used to deduplicate |
| Regular expressions | `preprocess/clean.go` | `regexp` collapses repeated whitespace |
| Value semantics | `preprocess/clean.go` | Structs passed by value; cleaning returns a fresh slice |

### A note on "inheritance"

Go has **no classical inheritance**. Instead it uses **composition** via struct embedding (a `Footballer` *has an* `Athlete`) and **interfaces** for shared behavior. This is more explicit than subclassing and avoids fragile deep class hierarchies.

### Data cleansing stages

- **Normalization** — trim/collapse whitespace and unify casing (`MESSI` → `Messi`).
- **Validation** — drop records missing required fields (e.g. empty name).
- **Sanitization** — fix out-of-range values (a negative goal count is clamped to `0`).
- **Deduplication** — collapse repeats using a composite `name|team` key.

---

## Sample Output

Your exact output will vary, but it follows this shape:

```
Raw: 4 footballers, 2 cricketers
Cleaned: 2 footballers, 1 cricketers

=== Cleaned Performers ===
Lionel Messi plays for Inter Miami (active) | 15 goals
Pele plays for Santos (retired) | 0 goals
Virat Kohli plays for Rcb (active) | 12000 runs
```

Notice that duplicates and the empty-name row were removed, casing was normalized, and Pele's negative goals were clamped to `0`.

---

## Next Steps

Ideas to extend the project as you keep learning:

- Make cleansing **configurable** with a `CleanOptions` struct (e.g. clamp vs. drop bad numbers).
- Factor the shared footballer/cricketer cleaning logic behind a **generic** pipeline.
- Read from **CSV** instead of JSON to practice another format.
- Expose the cleaned data over an **HTTP API**.
- Add **unit tests** (`go test`) for the cleansing functions.

---

## License

This is a learning project — use it freely.
