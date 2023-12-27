# Use an official Golang runtime as a parent image
FROM golang:1.21.5-bullseye

# Set the working directory inside the container
WORKDIR /restApi

# Copy the Go Modules manifests and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o GolandRestApi ./cmd/server

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./GolandRestApi"]

# Commands

# Build Image
# docker build -t golandrestapi .

# View a summary of image vulnerabilities and recommendations
# docker scout quickview

# Run the image in a container
# docker run -d -p 8080:8080 --env-file .env golandrestapi # So the variables inside the .env file are used in the container