FROM golang:1.22 AS build
WORKDIR /app

# Copy module file(s) first for better layer caching
COPY go.mod ./
# Copy the rest of the source
COPY . .

# Resolve dependencies and write go.sum inside the build
RUN go mod tidy

# Build the binary
RUN go build -o sports-app .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/sports-app .
COPY --from=build /app/players.json .
CMD ["./sports-app"]