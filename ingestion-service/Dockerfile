# Use the official Go image
FROM golang:alpine

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod ./

# Copy the source code
COPY . .

# Run the code directly (Bypasses binary format errors)
CMD ["go", "run", "./cmd/main.go"]