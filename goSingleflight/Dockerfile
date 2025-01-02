# Use an official Golang image as the base
FROM golang:1.23.4

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# # Run the tests
# CMD ["go", "test", "./..."]
