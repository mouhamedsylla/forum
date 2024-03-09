# Start from the Alpine Linux image with the latest version of Go installed
FROM golang:alpine

# Install sqlite3, bash, and necessary C libraries
RUN apk add --no-cache sqlite libgcc libstdc++ bash gcc musl-dev make

WORKDIR /app

# Copy your entire project
COPY . .

# Adjusted command to build the WASM file from the ./App directory
WORKDIR /app/App
RUN GOOS=js GOARCH=wasm go build -o ./internal/assets/main.wasm .

# Change back to the root project directory
WORKDIR /app

ENV CGO_ENABLED=1
# Build the main application executable
RUN go build -o forum ./cmd/main.go

# Set the working directory to /app for CMD execution
WORKDIR /app

RUN make app

WORKDIR /app

# Command to run the start script when the container starts
CMD ["./app"]