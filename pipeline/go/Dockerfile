# ---------- Build stage ----------
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod ./
COPY go.sum* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o minimarket-go-pipeline \
    ./main.go


# ---------- Runtime stage ----------
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/minimarket-go-pipeline .
COPY --from=builder /app/config ./config

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

CMD ["./minimarket-go-pipeline"]