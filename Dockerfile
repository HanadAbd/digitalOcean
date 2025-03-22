# Build stage
FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o myapp

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/myapp .
# Copy templates directory
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./myapp"]
