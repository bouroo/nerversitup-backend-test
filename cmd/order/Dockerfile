# Start from a Go image
FROM golang AS builder

# Set the working directory
WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download \
    && go mod verify

# Copy the application files
COPY . .

# Build the application
RUN go build -o order .

# Start from a clean debian image
FROM debian:stable-slim
LABEL maintainer="Kawin Viriyaprasopsook <kawin.v@kkumail.com>"

# Set timezone and locale
ARG	timezone="Asia/Bangkok"
ENV	LANG C.UTF-8
ENV	LC_ALL C.UTF-8
ENV	TZ $timezone

# Change locale
RUN echo $timezone > /etc/timezone \
    && cp /usr/share/zoneinfo/$timezone /etc/localtime

# Set the working directory
WORKDIR /app

# Copy the application binary from the builder image
COPY --from=builder /app/order .
COPY --from=builder /app/configs configs

# Expose the port that the application listens on
EXPOSE 8002

# Start the application
CMD ["./order"]
