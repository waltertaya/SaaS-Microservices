# Start from Go base image
FROM golang:1.23

# Set working directory
WORKDIR /app

# Copy go mod & sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build
RUN go build -o auth-service ./cmd/main.go

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./auth-service"]
