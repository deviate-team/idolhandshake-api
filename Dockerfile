# Use an official Golang runtime as a parent image
FROM golang:1.18-alpine AS build

# Set the working directory in the container to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the binary executable for the Fiber app
RUN go build -o main .

# Use a minimal Alpine Linux image as a parent image for the runtime
FROM alpine:3.14

# Set the working directory in the container to /app
WORKDIR /app

# Copy the binary executable from the build image to the runtime image
COPY --from=build /app/main .

# Expose the port that the Fiber app listens on
EXPOSE 3000

# Start the Fiber app
CMD ["./main"]