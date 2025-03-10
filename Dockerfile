# Build stage
FROM golang:1.24.1 AS builder
WORKDIR /app

# Copy go.mod and regenerate go.sum
COPY go.mod ./
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /proxy proxy.go

# Use Google Distroless (with CA store)
FROM gcr.io/distroless/static:nonroot 
WORKDIR /
COPY --from=builder /proxy /proxy

# Expose the port
EXPOSE 54878
USER nonroot 
CMD ["/proxy"]
