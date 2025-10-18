# Use the official Go image as a base image
FROM golang:1.25.3-alpine AS builder

# Add build arguments
ARG VERSION
ARG TARGETOS
ARG TARGETARCH

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application, injecting the version
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION}" -o main .

# Use a minimal base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the static folder
COPY --from=builder /app/static ./static/

# Copy the template file
COPY --from=builder /app/template.html .

# Copy the robots.txt file
COPY --from=builder /app/robots.txt .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
