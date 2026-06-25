FROM golang:1.22 AS build
WORKDIR /app

# Copy module file(s) first for better layer caching
COPY go.mod ./
# Copy the rest of the source
COPY . .

# Resolve dependencies and write go.sum inside the build
RUN go mod tidy

# Build the binary
RUN go build -o app .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/app .
COPY data.json .
CMD ["./app"]