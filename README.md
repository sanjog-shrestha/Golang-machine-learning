# go-data-demo

A beginner-friendly Go project demonstrating an end-to-end data workflow — **loading, cleansing, exploratory data analysis (EDA), and visualization** — running entirely in **Docker** via **Docker Compose**. No local Go installation required.

## Features

- **Load** JSON data from a file into typed Go structs
- **Cleanse & preprocess**: trim whitespace, normalize casing, validate fields, deduplicate records
- **EDA**: compute summary statistics (count, min, max, mean, median, standard deviation)
- **Visualize**: render an interactive HTML bar chart with [go-echarts](https://github.com/go-echarts/go-echarts)
- Fully containerized — runs with a single `docker compose up`

## Project Structure

```
go-data-demo/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── data.json
├── main.go
└── out/              # generated: chart.html lands here
```

## Requirements

- [Docker](https://www.docker.com/) with the Compose plugin

That's it — Go runs inside the container, so it does not need to be installed on the host.

## Getting Started

Create the output folder and run:

```bash
mkdir -p out
docker compose up --build
```

This builds the image (resolving dependencies and compiling inside the container), runs the program, and writes `chart.html` into `./out/`. Open `out/chart.html` in a browser to view the interactive chart.

### Common Commands

| Command | What it does |
|---------|--------------|
| `docker compose up --build` | Rebuild and run (use after code changes) |
| `docker compose up` | Run again without rebuilding |
| `docker compose down` | Stop and clean up |

## The Data Pipeline

The program processes records through a deliberate sequence.

### 1. Load

Reads `data.json` into memory and unmarshals it into a slice of `Person` structs.

```go
type Person struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}
```

### 2. Cleanse & Preprocess

Each record is normalized, then validated, then deduplicated. Order matters — normalizing before validating prevents, for example, a whitespace-padded email from being wrongly rejected.

| Step | Rule |
|------|------|
| Clean | Trim whitespace, lowercase emails, capitalize names |
| Validate | Drop records with empty name, age outside 0–120, or no `@` in email |
| Deduplicate | Drop repeated IDs (keeps the first valid one) |

### 3. EDA

Computes summary statistics over the age column.

| Statistic | Meaning |
|-----------|---------|
| Count | Number of valid records |
| Min / Max | Smallest and largest values |
| Mean | Average |
| Median | Middle value |
| Std Dev | Spread around the mean (population) |

### 4. Visualize

Renders an interactive bar chart of age per person to `chart.html` using go-echarts (HTML/JavaScript output, so no GUI or image libraries are needed in the container).

## Sample Output

```
Loaded 6 raw records
Kept 3 clean records, dropped 3

--- EDA: Age ---
Count:   3
Min:     25.00
Max:     35.00
Mean:    30.00
Median:  30.00
Std Dev: 4.08

Wrote chart.html
```

## Go Concepts Demonstrated

- **Struct tags** map JSON keys to Go fields
- **Pointer vs value receivers**: `clean()` uses a pointer receiver (it mutates), `valid()` uses a value receiver (it only reads)
- **Maps as sets**: `map[int]bool` tracks seen IDs, since Go has no built-in set type
- **Slice copying before sorting**: `sort` mutates in place, so a copy protects the caller's data
- **Multi-stage Docker builds** keep the final image small

## Next Steps

- Add a histogram to show age distribution by bucket
- Compute frequency analysis on categorical fields
- Read CSV input with `encoding/csv`
- Stream large files instead of loading them fully into memory
- Pull data from a database with `database/sql`

## License

MIT
