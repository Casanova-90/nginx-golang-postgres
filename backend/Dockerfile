# Use Go base image that matches your version (e.g., Go 1.15 for compatibility)
FROM golang:1.15-alpine

WORKDIR /app

# Copy the application source code
COPY . .

# Install dependencies (if you have any)
RUN go mod tidy

# Build the Go application
RUN go build -o server .

# Use a lightweight image for running the application
FROM alpine:3.12

WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=0 /app/server .

# Expose the port on which the app will run
EXPOSE 8000

# Start the application
CMD ["./server"]
