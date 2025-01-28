# Use the official Golang image as the build stage
FROM golang:1.23.5-alpine3.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy all files from the current directory to the working directory inside the container
COPY . .

# Build the Go application and output the binary as 'main'
RUN go build -o main main.go

# Use the official Alpine image as the run stage
FROM alpine:3.20

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the build stage to the current stage
COPY --from=builder /app/main .
COPY start.sh .
COPY app.env .
COPY wait-for.sh .
COPY db/migration ./db/migration

# Expose port 9091 to the outside world
EXPOSE 9091

# Command to run the application
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]