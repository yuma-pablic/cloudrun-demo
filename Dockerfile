FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/main /usr/local/bin/main

CMD ["main"]