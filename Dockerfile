# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o test-assigment-ltp

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/test-assigment-ltp .
EXPOSE 8080
CMD ["./test-assigment-ltp"]