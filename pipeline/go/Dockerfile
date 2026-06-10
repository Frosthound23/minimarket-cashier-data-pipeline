FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o minimarket-go-pipeline ./main.go


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/minimarket-go-pipeline .
COPY config ./config

CMD ["./minimarket-go-pipeline"]