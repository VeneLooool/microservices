FROM golang:1.18.2-alpine3.15 as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o warehouse ./cmd/warehouse

EXPOSE 8080

ENTRYPOINT ["/app/warehouse"]