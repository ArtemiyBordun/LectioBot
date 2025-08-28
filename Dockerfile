FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/internal/config/config.env ./internal/config/config.env
COPY --from=builder /app/internal/config/course_config.json ./internal/config/course_config.json
COPY --from=builder /app/internal/config/credentials.json ./internal/config/credentials.json

CMD ["./main"]
