# Stage 1: Build
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod ./
# COPY go.sum ./ (If you have it)
COPY . .
RUN go build -o main ./cmd/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8081
CMD ["./main"]