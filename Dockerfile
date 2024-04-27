FROM golang:1.22-alpine AS builder

ARG VER=dev

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

RUN apk add --update make git gcc libc-dev curl file

# Build the Go binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    export CGO_ENABLED=1 &&\
    export LDFLAGS="-X main.version=$VER -linkmode external -extldflags \"-static\"" &&\
    go build -v -ldflags "${LDFLAGS}" -o binomeme ./cmd/binomeme/main.go &&\
    file binomeme

# Create a minimal production image
FROM alpine:latest

# It's essential to regularly update the packages within the image to include security patches
# Reduce image size
RUN apk update && \
    apk upgrade && \
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

# Avoid running code as a root user
RUN adduser -D binomeme
USER binomeme

# Set the working directory inside the container
WORKDIR /app

RUN mkdir /app/data

# Copy only the necessary files from the builder stage
COPY --from=builder /app/binomeme .

# Run the binary when the container starts
CMD ["./binomeme"]
