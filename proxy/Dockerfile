FROM golang:1.20.3-bullseye AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o main

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]