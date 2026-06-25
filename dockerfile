# Build Stage
FROM golang:1.22 AS build 
WORKDIR /app
COPY go.mod ./
COPY . . 
RUN go build -o app .

# Run stage
FROM debian:bookworm-slim
WORKDIR /app 
COPY --from=build /app/app .
COPY data.json .
CMD [ "./app" ]