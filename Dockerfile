FROM golang:1.25 AS builder

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api/main.go

# --- runtime stage ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]