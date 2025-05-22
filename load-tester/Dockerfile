# Dockerfile
FROM golang:1.23

WORKDIR /app

# RUN apt-get update && \
#     apt-get install -y protobuf-compiler


# RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
#     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# COPY proto/common.proto .

# Generate Go code (protoc plugins are not needed since we use pre-generated go.mod)
# RUN protoc --go_out=. --go-grpc_out=. greeter.proto
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main

CMD ["./main"]
