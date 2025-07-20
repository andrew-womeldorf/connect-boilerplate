# Connect RPC API Boilerplate

This is a boilerplate repository for building Go APIs using [Connect
RPC](https://connectrpc.com/) and [Buf](https://buf.build/). It demonstrates a
protocol-first development approach where all interactions with the system are
driven through protobuf service definitions.

## Architecture Philosophy

**Schema-First Development**: This boilerplate enforces that all access to
business logic and data must flow through the protobuf-defined service
interface. There are no direct paths to business logic or database operations
outside of what's explicitly defined in the `.proto` files.

This approach ensures:
- **API Contract Clarity**: The protobuf definitions serve as the single source
  of truth for what the API can do
- **Type Safety**: All data structures and service methods are strongly typed
- **Multi-Language Support**: Connect RPC works with any language that supports
  protobuf
- **Backwards Compatibility**: Changes to the API can be validated for breaking
  changes using buf

## Getting Started

### First-Time Setup (Required)

**IMPORTANT**: If you've copied or forked this repository, run the setup script
**once** to rename the project:

```bash
./setup.sh your.domain.com/username/project-name
```

For example:
```bash
./setup.sh github.com/myuser/my-api
./setup.sh git.company.com/team/project-api
```

This script will:
- Update all Go imports from `github.com/andrew-womeldorf/connect-boilerplate`
  to your module name
- Update all configuration files (go.mod, buf.gen.yaml, etc.)

After running setup:
1. **Generate protobuf code**: `mise run proto:generate`
2. **Run quality checks**: `mise run check`
3. **Commit your changes**: `git add . && git commit -m "Rename project to
   your-module-name"`
4. **Remove setup script**: `rm setup.sh`

### Development Workflow

1. **Trust mise configuration**: `mise trust`
2. **Generate protobuf code**: `mise run proto:generate`
3. **Start the server**: `mise run serve`
4. **Test with RPC commands**: `./build/api user list-users` (after building)

## Directory Structure

### `proto/`
**The Source of Truth** - Contains all protobuf definitions using the 1-1-1 pattern:
- `proto/user/v1/user_service.proto` - Service definition (imports all message types)
- `proto/user/v1/user.proto` - Core entity definitions
- `proto/user/v1/{operation}.proto` - Individual request/response message pairs

**Key Principle**: All functionality must be defined here first. No business
logic should exist without a corresponding protobuf definition.

### `gen/`
**Generated Code** - Auto-generated Go code from protobuf definitions. Never
edit manually.
- `gen/user/v1/*.pb.go` - Protobuf message types
- `gen/user/v1/userv1connect/*.connect.go` - Connect RPC service interfaces

### `internal/services/user/`
**Business Logic Layer** - Implements the actual service logic defined in protobuf.
- `service.go` - Core business logic that implements the protobuf-generated interfaces
- Methods must match exactly what's defined in the `.proto` service definitions

### `internal/server/`
**HTTP Server Layer** - Bridges the business logic to HTTP transport and handles data access.
- `server.go` - HTTP server setup with Connect RPC handlers, gRPC reflection, and h2c support
- `user_connect_handler.go` - Thin adapter that connects
  `internal/services/user` service to Connect RPC interface
- Provides both standalone server (`Run()`) and handler creation (`CreateHandler()`) for Lambda
- **Database Access**: All data persistence should be handled at this layer

### `cmd/`
**Entry Points** - Two distinct deployment targets:

#### `cmd/cli/`
**CLI Interface** - Cobra-based command-line tool with dual-mode operation:
- `serve.go` - Starts the HTTP server
- `user.go` - User RPC client setup with endpoint configuration
- `user_*.go` - Individual User RPC command implementations (one file per RPC method)

**Dual Mode Support**:
- **In-memory mode** (default): Directly calls service methods for testing/development
- **Remote mode** (`--endpoint` flag): Makes HTTP Connect RPC calls to running server

#### `cmd/lambda/`
**AWS Lambda Entry Point** - Uses the same server handler for serverless deployment.

### Configuration Files

- `buf.yaml` & `buf.gen.yaml` - Buf configuration for protobuf linting and code generation
- `mise.toml` - Task runner configuration for development workflows
- `go.mod` - Go module with Connect RPC and related dependencies

## Development Workflow

1. **Define the API**: Start by defining or modifying `.proto` files in `proto/`
2. **Generate Code**: Run `mise run proto:generate` to update generated code
3. **Implement Business Logic**: Add implementation in
   `internal/services/<svc>/service.go`
4. **Update Adapter**: Ensure `internal/server/<svc>_connect_handler.go` routes
   to your service methods
5. **Test**: Use CLI commands or start server to test functionality

## Key Commands

- `mise run proto:generate` - Regenerate code from protobuf definitions
- `mise run serve` - Start development server on port 8088
- `mise run check` - Run all linters, formatters, and tests
- `mise run build` - Build CLI binary to `build/api`

## Protocol-First Enforcement

This boilerplate is designed to prevent common anti-patterns:

- ❌ **Don't**: Create RPC methods that bypass the protobuf service definitions
- ❌ **Don't**: Add database access outside of the `pkg/api` layer  
- ❌ **Don't**: Create business logic that isn't exposed through a protobuf RPC method

- ✅ **Do**: Define all functionality in `.proto` files first
- ✅ **Do**: Implement business logic in `pkg/api/service.go` methods
- ✅ **Do**: Use the CLI RPC commands for testing and client interaction

This ensures your API remains consistent, type-safe, and evolvable while
supporting multiple deployment targets (HTTP server, Lambda, CLI).
