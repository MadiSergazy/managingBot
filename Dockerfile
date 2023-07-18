# Use the official Golang image as the base image
FROM golang:1.20-alpine3.18
# FROM gol
# Use the official Gang:1.20-alpine3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files to the working directory
COPY go.mod go.sum ./

# Download and cache Go modules dependencies
RUN go mod download

# Copy the rest of the project files to the working directory
COPY . .

# Install MySQL client
RUN apk update && apk add mysql-client

# Build the Go application
RUN go build -o main ./cmd/main.go

# Expose the port on which the application will run
# EXPOSE 8080

# Set the entry point command for the container
CMD ["./main"]
# COPY --from=0 it means we want to copy from previous step of building
# EXPOSE 80 - port for getting application in inside

# ENV DB_NAME=telegrambot
# ENV DB_HOST=127.0.0.1 
# ENV DB_PORT=3306
# ENV DB_USER=Lift_kz
# ENV DB_PASSWORD=Lift@2023

# docker build -t my-golang-app .
# docker run  my-golang-app