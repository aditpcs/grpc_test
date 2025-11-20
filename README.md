# gRPC Test Client

This is a sample gRPC client for testing the PgCall service from the 2507_pose-notifications-middleware project.

## Prerequisites

- Go 1.23.0 or higher
- The gRPC server running on `localhost:50051`
- protoc (Protocol Buffer Compiler) - if you need to regenerate proto files

## Project Structure

```
grpc_tets/
├── main.go              # Simple gRPC client example
├── client_example.go    # Advanced client with CLI flags
├── bridge_http/         # Generated gRPC code
│   ├── bridge_http.proto
│   ├── bridge_http.pb.go
│   └── bridge_http_grpc.pb.go
├── go.mod
├── go.sum
└── README.md
```

## Setup

1. Navigate to the grpc_tets directory:
   ```bash
   cd ~/go/src/grpc_tets
   ```

2. Download dependencies:
   ```bash
   go mod tidy
   ```

3. Build the clients:
   ```bash
   go build -o grpc_client main.go
   go build -o grpc_client_advanced client_example.go
   ```

## Usage

### Simple Client (main.go)

Run the basic client with hardcoded values:
```bash
go run main.go
```

Or use the compiled binary:
```bash
./grpc_client
```

### Advanced Client (client_example.go)

Run with default values:
```bash
go run client_example.go
```

Run with custom parameters:
```bash
go run client_example.go \
  -server localhost:50051 \
  -vendor my-vendor \
  -username my-user \
  -password my-password \
  -reqtype POST \
  -timeout 30s
```

Or use the compiled binary:
```bash
./grpc_client_advanced -vendor test-vendor -username testuser
```

### Available Flags (client_example.go)

- `-server`: gRPC server address (default: `localhost:50051`)
- `-vendor`: Vendor name (default: `example-vendor`)
- `-username`: Username (default: `example-user`)
- `-password`: Password (default: `example-password`)
- `-reqtype`: Request type (default: `GET`)
- `-timeout`: Request timeout duration (default: `10s`)

## Regenerating Proto Files

If you modify the `.proto` file, regenerate the Go code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       bridge_http/bridge_http.proto
```

## Example Output

```
2024/11/17 15:51:37 Connecting to gRPC server at localhost:50051...
2024/11/17 15:51:37 Connected successfully!
2024/11/17 15:51:37 ---
2024/11/17 15:51:37 Sending request:
2024/11/17 15:51:37   Vendor:   example-vendor
2024/11/17 15:51:37   Username: example-user
2024/11/17 15:51:37   ReqType:  GET
2024/11/17 15:51:37 ---
2024/11/17 15:51:37 Response received:
2024/11/17 15:51:37   Status Code: 200
2024/11/17 15:51:37   Response Body Length: 1024 bytes
2024/11/17 15:51:37   Headers:
2024/11/17 15:51:37     Content-Type: application/json
2024/11/17 15:51:37   Response Body: {"status":"success"}
2024/11/17 15:51:37 ---
2024/11/17 15:51:37 Request completed successfully!
```

## Notes

- Make sure the gRPC server is running before executing the client
- The client uses insecure credentials for testing purposes
- Modify the server address if your gRPC server is running on a different host/port
