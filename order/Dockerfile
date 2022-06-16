FROM golang:1.18.2-alpine3.15 as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o orders ./cmd/orders

EXPOSE 8082

ENTRYPOINT ["/app/orders"]