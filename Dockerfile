## Dockerfile
FROM golang:1.21 as builder
WORKDIR /app

# Copy source files
COPY go.mod ./
RUN go mod download
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /proxy proxy.go

# Create minimal scratch image
FROM scratch
COPY --from=builder /proxy /proxy

# Expose the port
EXPOSE 54878
CMD ["/proxy"]
