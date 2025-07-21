# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

All commands are managed through mise. The project requires `mise trust` before first use.

### Development Workflow
- `mise run proto:generate` - Generate Go code from protobuf definitions (required after proto changes)
- `mise run serve` - Run the API server locally on port 8080
- `mise run cli` - Invoke the CLI
- `mise run check` - Run all formatters, linters, and tests
- `mise run build` - Build the CLI binary to `build/api`

### Testing and Quality
- `mise run test` - Run all tests
- `mise run test -- -run TestSpecificFunction` - Run a specific test
- `mise run lint` - Run golangci-lint
- `mise run format` - Format code (go fmt, go mod tidy, goimports)

### Protobuf Management
- `mise run proto:lint` - Lint protobuf files
- `mise run proto:generate` - Generate Go code from protobuf definitions (required after proto changes)

## Architecture

This is a Connect RPC API (gRPC-compatible) with the following structure:

### Protocol Buffers (1-1-1 Pattern)
- `proto/user/v1/` - Uses the 1-1-1 pattern: one proto file per message type
- `user_service.proto` - Service definition only (imports message files)
- Individual message files: `list_users.proto`, `get_user.proto`, etc.
- Generated code outputs to `gen/user/v1/`

### Core Components
1. **Service Layer** (`internal/services/<service>/service.go`)
   - Business logic implementation
   - Methods return Connect objects directly
   - Each service method has its own file (`op_*.go`)
   - Each service has its own store interface and implementation(s) (`internal/services/<service>/store/*`)

2. **Server Layer** (`internal/server/`)
   - `server.go` - HTTP server setup with Connect RPC handlers, gRPC reflection, and h2c support
   - `<service>_connect_handler.go` - Adapts the service to the Connect interface
   - Provides both `Run()` for standalone server and `CreateHandler()` for Lambda

3. **CLI** (`cmd/cli/`)
   - Cobra-based with three main command groups:
     - `serve` - Start the API server
     - `user` - Client commands for the User service with dual mode support:
       - In-memory mode (default): Direct service calls
       - Remote mode (`--endpoint` flag): HTTP client calls
   - Each User RPC method has its own file (`user_*.go`)

4. **Lambda** (`cmd/lambda/`)
   - AWS Lambda entry point using the server's `CreateHandler()`

### Client Usage Patterns
The CLI RPC commands support two modes:
- **In-memory**: `./api rpc list-users` (creates service directly)
- **Remote**: `./api rpc list-users --endpoint http://localhost:8080` (HTTP client)

All RPC commands output JSON and support structured flags for input parameters.
