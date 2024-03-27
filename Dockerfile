# Start with a base Go image to build your application
FROM golang:1.18 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stackScript stackScript.go

# Use a small base image
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/stackScript .
COPY --from=builder /app/heapScript .

# Expose the port the app runs on
EXPOSE 6060

# Command to run the binary (this can be overridden)
CMD ["./stackScript"]  # Default command, can be overridden when running the container
