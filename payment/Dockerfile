FROM golang:1.18.2-alpine3.15 as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o payment ./cmd/payment

EXPOSE 8081

ENTRYPOINT ["/app/payment"]