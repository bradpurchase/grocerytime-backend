FROM golang:alpine AS builder

# Set env variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk update && apk add --no-cache git

# Move to working directory /build
WORKDIR /build

# Copy and download dependencies using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy code to the container
COPY . .

# Build the go app
RUN go build -o main .

# Stage 2
FROM scratch 

WORKDIR /server

COPY --from=builder /build/main .
COPY --from=builder /build/.env .

EXPOSE 8080

CMD ["./main"]