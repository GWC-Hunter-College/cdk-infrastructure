# Edited from https://docs.aws.amazon.com/lambda/latest/dg/go-image.html, 
# 'Using an AWS OS-only base image -> creating an image from the provided.al2023 base image' section #4.
FROM golang:1.24.4-alpine3.21 AS build
WORKDIR /init

# Copy dependencies list
COPY go.mod go.sum ./

# Build
COPY main.go .
RUN go build -o main main.go

# Copy artifacts to a clean image
FROM alpine:3.21
COPY --from=build /init/main /main
ENTRYPOINT [ "/main" ]