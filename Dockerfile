# Stage 1: Build the binary
FROM golang:1.24-bookworm AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy project source code
COPY . .

# Receive the token as build-arg
ARG GITHUB_TOKEN

# Define which modules are private
ENV GOPRIVATE=github.com/{{{YOUR_PROJECT}}}

# Configure git to use the token and run go mod tidy
RUN git config --global url."https://$GITHUB_TOKEN@github.com/".insteadOf "https://github.com/" && go mod tidy

# Download dependencies
RUN go mod tidy

# Generate swagger docs
RUN go tool swag init -g ./api/main.go

# Build the binary
RUN go build -o main .

# Stage 2: Create the final image
FROM debian:bookworm-slim

# Install required dependencies
RUN apt-get update && apt-get install -y ca-certificates

# Copy the generated binary from the build stage
COPY --from=builder /app/main /

# Command to start the service
CMD ["/main"]
