# # Stage 1: Build the application
# FROM golang:1.20-alpine3.18 as builder

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the Go modules files to the working directory
# COPY go.mod go.sum ./

# # Download and cache Go modules dependencies
# RUN go mod download

# # Copy the rest of the project files to the working directory
# COPY . .

# # Install Vim and the necessary packages for setting up locales
# RUN apk update && apk add vim musl-locales

# # Set up the Russian locale
# ENV LANG=ru_RU.utf8
# ENV LANGUAGE=ru_RU.utf8
# ENV LC_ALL=ru_RU.utf8

# # Build the Go application
# RUN go build -o main ./cmd/main.go

# # Stage 2: Create the final production image
# FROM alpine:3.14 as production

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the binary from the builder stage to the final image
# COPY --from=builder /app/main .
# COPY --from=builder /app/.env .

# # Set the entry point command for the container
# CMD ["./main"]






# it was last
# Stage 1: Build the application
FROM golang:1.20-alpine3.18 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files to the working directory
COPY go.mod go.sum ./

# Download and cache Go modules dependencies
RUN go mod download

# Copy the rest of the project files to the working directory
COPY . .


# Install MySQL client
# RUN apk update && apk add mysql-client

# Build the Go application
RUN go build -o main ./cmd/main.go

# Stage 2: Create the final production image
FROM alpine:3.14 as production

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Set the entry point command for the container
CMD ["./main"]

# ENV DB_NAME=telegrambot
# ENV DB_HOST=127.0.0.1 
# ENV DB_PORT=3306
# ENV DB_USER=Lift_kz
# ENV DB_PASSWORD=Lift@2023

# docker build -t my-golang-app .
# docker run  my-golang-app


# # Stage 1: Build the application
# FROM golang:1.17 as builder

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the Go modules files to the working directory
# COPY go.mod go.sum ./

# # Download and cache Go modules dependencies
# RUN go mod download

# # Copy the rest of the project files to the working directory
# COPY . .

# # Build the Go application
# RUN go build -o main ./cmd/main.go

# # Stage 2: Create the final production image
# FROM ubuntu:20.04 as production

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the binary from the builder stage to the final image
# COPY --from=builder /app/main .
# COPY --from=builder /app/.env .

# # Set the entry point command for the container
# CMD ["./main"]

# # Set the LANG environment variable to en_US.UTF-8
# ENV LANG en_US.UTF-8
