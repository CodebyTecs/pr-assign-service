FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pr-assign-service ./cmd/pr-assign-service

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/pr-assign-service /app/pr-assign-service

EXPOSE 8080

ENV HTTP_SERVER_ADDRESS=0.0.0.0
ENV HTTP_SERVER_PORT=8080
ENV HTTP_SERVER_TIMEOUT=4s

CMD ["/app/pr-assign-service"]