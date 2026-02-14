# Build stage (docs/ is not copied – generated inside image from handler source)
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -d . -o docs && \
	CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# Run stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /api .
RUN mkdir -p uploads
EXPOSE 8080
CMD ["./api"]
