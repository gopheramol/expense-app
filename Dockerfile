# Start from the official Go image as the base
FROM golang:1.22-alpine AS build

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the current working directory
COPY . .

# Build the Go application
RUN go build -o expense-tracker

# Start a new stage from scratch
FROM alpine:latest  

# Set the working directory to /app in the new stage
WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=build /app/expense-tracker .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./expense-tracker"]
