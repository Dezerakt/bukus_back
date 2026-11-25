# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/app

# Final stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app

COPY --from=builder /bin/app /bin/app
COPY --from=builder /app/migrations /migrations

USER app

ENTRYPOINT ["/bin/app"]
