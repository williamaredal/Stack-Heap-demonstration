# Use a minimal Go image as the base
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o myapp .

# Use a minimal base image for the runtime container
FROM alpine:latest

# Set the working directory inside the runtime container
WORKDIR /app

# Copy the built binary from the build container
COPY --from=build /app/myapp .

# Command to run the Go application
CMD ["./myapp"]

