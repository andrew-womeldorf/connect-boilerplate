package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1/examplev1connect"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/server"
)

var (
	apiEndpoint string
)

// rpcCmd represents the rpc command
var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Execute RPC calls to the User service",
	Long: `Execute RPC calls to the User service using RPC-style commands.
This command provides subcommands for all RPCs in the User service.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		if err := cmd.Help(); err != nil {
			slog.Error("Failed to display help", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(rpcCmd)

	// Add API endpoint flag to the rpc command
	rpcCmd.PersistentFlags().StringVar(&apiEndpoint, "endpoint", "", "API endpoint URL (e.g., http://localhost:8088)")

	// Add all RPC commands
	rpcCmd.AddCommand(listUsersCmd())
	rpcCmd.AddCommand(getUserCmd())
	rpcCmd.AddCommand(createUserCmd())
	rpcCmd.AddCommand(updateUserCmd())
	rpcCmd.AddCommand(deleteUserCmd())
}

// getClient returns either a local client or a remote client based on whether the API endpoint is provided
func getClient(ctx context.Context) (examplev1connect.UserServiceClient, error) {
	if apiEndpoint != "" {
		// Use Connect client with remote endpoint
		httpClient := http.DefaultClient
		return examplev1connect.NewUserServiceClient(
			httpClient,
			apiEndpoint,
		), nil
	} else {
		// Use local service with ServiceAdapter
		srv := server.NewServer(0) // Port doesn't matter for local client
		service, err := srv.GetService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get service: %w", err)
		}
		return server.NewServiceAdapter(service), nil
	}
}

// printJSON prints the given data as JSON
func printJSON(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal JSON", "error", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonBytes))
}
