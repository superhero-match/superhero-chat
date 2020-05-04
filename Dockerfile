# Start from the latest golang base image.
FROM golang:latest as builder

# Add Maintainer Info.
LABEL maintainer="Nikolajus Karpovas <mwsoftofficial@gmail.com>"

# Create build directory.
RUN mkdir build

# Copy project to build directory.
COPY . /build

# Set build as working directory.
WORKDIR /build/cmd/chat

# Fetch dependencies.
RUN go mod download

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main

# Set build as working directory.
WORKDIR /build/cmd/health

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o health

# Create unprivelleged user.
RUN adduser --disabled-login appuser

# Start a new stage from scratch.
FROM alpine:latest

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

# Create build directory.
RUN mkdir app

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /build/cmd/chat/main /app/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /build/cmd/health/health /app/

# Copy the config file from the previous stage.
COPY ./config.yml /app/

# Copy the certificates from the previous stage.
COPY ./certificate.pem /app/cmd/chat/
COPY ./key.pem /app/cmd/chat/

# Set working directory in current stage.
WORKDIR /app

# Use an unprivileged user.
USER appuser

# Expose port 8120.
EXPOSE 8120 8120

# Command to run the executables.
CMD ["sh", "-c", "( ./health & ) && ./main"]