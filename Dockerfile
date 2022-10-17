# syntax=docker/dockerfile:1

# Google already build an image for Go, just inherit it for use.
FROM golang:latest

# This working directory is referencing inside the image, so this is fine.
WORKDIR /Project

# Copy where my current system's directory's stuff into the image
COPY . .

# Pass in my environment variables
ENV Loggly_Token=http://logs-01.loggly.com/inputs/5e085983-7ed1-4fc1-bf95-5f6278278035/tag/http/

# Tell the image to build and run the Go code
RUN go build main.go
CMD [ "go", "run", "main.go" ]