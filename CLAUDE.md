## Bash Commands

- `mise run format` - run all auto formatters
- `mise run lint` - runs all linters
- `mise run test` - run all tests
- `mise run test -- -run TestSpecificFunction` - run specific test
- `mise run check` - runs all formatters, linters, and tests 
- `mise run build` - builds the cli to `build/api`
- `mise run cli` - builds and runs `./cmd/cli`

## Workflow

- A task is not finished until `mise run check` completes with a zero exit code.

## Structure

- The API Interface is defined in `proto/<service>/<version>`. Changes to the interface must occur here first.
- All rpc methods must be handled by:
    - `internal/server/<service>_connect_handler.go`, which unwraps ConnectRPC Requests, calls the service, and creates new ConnectRPC Responses.
    - `internal/services/<service>/op_<rpc>.go`, which handles the business logic for a `pb.<rpc>Request`.
    - `cmd/cli/<service>/op_<rpc>.go`, which adds a new cobra subcommand with flags to support the Request object.
- The web interface (`internal/web/`) can be implemented in any way. But it must make use of `internal/services/*` for data management. It is not allowed to access the data store directly.
- Each service maintains its own store (`internal/services/<service>/store`)
    - The store interface is maintained in `internal/services/<service>/store/store.go`
    - The store is tested in `internal/services/<service>/store/store_test.go`, which must test all store implementations with the same inputs and outputs.
    - The store is implemented for both sqlite and dynamodb in `internal/services/<service>/store/<impl>/`
