# Step 1: Use the official Golang image to build the app
# This uses a specific Go version for stability
FROM golang:1.23-alpine as builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Install required dependencies for Go
RUN apk add --no-cache git

# Step 4: Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Step 5: Download all Go module dependencies
# Dependencies will be cached if the go.mod and go.sum files haven't changed
RUN go mod tidy

# Step 6: Copy the rest of the application code to the container
COPY . .

# Step 7: Build the Go app
RUN go build -o main .

# Step 8: Use a minimal base image for running the app
FROM alpine:3.20

# Step 9: Set the working directory for the minimal base image
WORKDIR /app

# Step 10: Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Step 11: Copy the Go app from the builder stage
COPY --from=builder /app/main .

# Step 12: Copy Configuration
COPY config-prod.toml .

# Step 13: Command to run the app
CMD ["./main", "-c", "config-prod.toml"]