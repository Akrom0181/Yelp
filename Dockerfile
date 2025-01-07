# # Step 1: Build the Go binary
# FROM golang:1.23 AS builder

# # Set the working directory inside the container
# WORKDIR /app

# # Copy go.mod and go.sum files first (for caching dependencies)
# COPY go.mod go.sum ./

# # Download and cache Go modules
# RUN go mod download

# # Copy the rest of the application code
# COPY . .

# # Build the Go binary
# RUN GOOS=linux GOARCH=amd64 go build -tags migrate -o yelpbinary ./cmd/app

# # Step 2: Create a lightweight image to run the binary
# FROM alpine:latest

# # Set the working directory inside the container
# WORKDIR /app

# # # Copy the binary from the builder stage
# COPY --from=builder /app/yelpbinary /app/
# COPY --from=builder /app/config /app/config
# COPY --from=builder /app/migrations /app/migrations

# # Expose the port your application uses (optional)
# EXPOSE 8080

# # Command to run the binary
# CMD ["/app/yelpbinary"]

# Step 1: Build the Go binary
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Debug: Check contents of ./cmd/app
RUN echo "Checking ./cmd/app contents:" && ls -l /app/cmd/app

RUN GOOS=linux GOARCH=amd64 go build -tags migrate -o yelpbinary ./cmd/app

# Debug: Verify binary exists
RUN echo "Checking binary:" && ls -l /app/yelpbinary

# Step 2: Create a lightweight image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/yelpbinary /app/
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

RUN apk add --no-cache libc6-compat

EXPOSE 8080

CMD ["/app/yelpbinary"]
