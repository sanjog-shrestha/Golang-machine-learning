# go-data-demo

A minimal Go project demonstrating how to **load and handle data in Go**, packaged to run in **Docker**. It reads JSON from a file into Go structs, processes and filters the data, and serializes the result back to JSON.

## Features

- Reads a file into memory with `os.ReadFile`
- Parses JSON into typed structs using `encoding/json`
- Maps JSON keys to Go fields via struct tags
- Filters a slice and serializes results back to JSON
- Multi-stage Docker build for a small final image

## Project Structure

```
go-data-demo/
├── Dockerfile
├── go.mod
├── data.json
└── main.go
```

## Requirements

- [Docker](https://www.docker.com/) (no local Go install needed), **or**
- [Go 1.22+](https://go.dev/dl/) if running natively

## Getting Started

### Run with Docker

Build the image:

```bash
docker build -t go-data-demo .
```

Run the container:

```bash
docker run --rm go-data-demo
```

### Run natively

```bash
go run main.go
```

### Develop without rebuilding

Mount your code into the official Go image and run it directly. Useful while iterating:

```bash
docker run --rm -v "$PWD":/app -w /app golang:1.22 go run main.go
```

## How It Works

The program loads `data.json`, unmarshals it into a slice of `Person` structs, prints each record, filters for people over 28, and prints the filtered set as indented JSON.

```go
type Person struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

| Concept        | What it does                                  |
|----------------|-----------------------------------------------|
| `os.ReadFile`  | Loads a file into a byte slice                |
| `json.Unmarshal` | Parses JSON bytes into Go values            |
| Struct tags    | Map JSON keys to struct fields                |
| `append`       | Builds dynamic slices for filtered results    |
| `json.MarshalIndent` | Serializes Go values to formatted JSON  |

## Sample Output

```
1: Alice (30)
2: Bob (25)
3: Carol (35)

Over 28: 2 people
[
  {
    "id": 1,
    "name": "Alice",
    "age": 30
  },
  {
    "id": 3,
    "name": "Carol",
    "age": 35
  }
]
```

## Next Steps

- Read CSV files with `encoding/csv`
- Stream large files instead of loading them fully into memory
- Pull data from a database with `database/sql`

## License

MIT
