# Use the latest official Golang image to build the application
FROM golang:latest AS builder

# Build arguments
ARG TARGETOS
ARG TARGETARCH

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the template file
COPY --from=builder /app/template.html .

# Copy the static file
COPY --from=builder /app/static ./static

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
