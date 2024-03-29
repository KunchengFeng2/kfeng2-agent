# Start the Go app build
FROM golang:latest AS build

# Copy source
WORKDIR /csc482/kfeng2-agent/source
COPY . .

# Get required modules
RUN go mod tidy

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# New build phase -- create binary-only image
FROM alpine:latest

# Add support for HTTPS
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates

WORKDIR /csc482/kfeng2-agent

# Copy files from previous build container
COPY --from=build /csc482/kfeng2-agent/source/main ./

# Add environment variables
# They will be defined in AWS ECS Task Definitions

# Check results
RUN env && pwd && find .

# Start the application
CMD ["./main"]

# Remember start the image with
# docker run -p <local port>:<inside port> -d <image name>