# Build the Go app
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app main.go

# Run the Go app 
FROM alpine AS runtime 
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["app"]