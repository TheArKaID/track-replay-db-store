FROM golang:alpine

# Package Installer
RUN apk update && apk add --no-cache git

# From this directory
WORKDIR /src

# Copy all files from current directory to /app in container
COPY . .

# Download all dependencies
RUN go mod tidy

# Build the Go app
RUN go build -o main

# From this directory
WORKDIR /app

# Copy main to /app
RUN cp /src/main /app/main

# delete all files in /src
RUN rm -rf /src

# GO!
ENTRYPOINT ["/app/main"]