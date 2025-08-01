FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o notepad ./cmd/app/app.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/* 

WORKDIR /app

COPY --from=builder /app/notepad .

EXPOSE 8000

CMD ["./notepad"]