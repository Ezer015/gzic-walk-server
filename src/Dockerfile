# Use the official Golang image as the base image
FROM golang:1.23.2-alpine

# Set Go proxy
ENV GOPROXY=https://goproxy.cn,direct

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy project files
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the application runs on
EXPOSE 8080

# Run the application
CMD ["./main"]