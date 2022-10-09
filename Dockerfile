# syntax=docker/dockerfile:1

# Google already build an image for Go, just inherit it for use.
FROM golang:latest

# This working directory is referencing inside the image, so this is fine.
WORKDIR /Project

# Copy where my current system's directory's stuff into the image
COPY . .

# Tell the image to build and run the Go code
RUN go build main.go
CMD [ "go", "run", "main.go" ]