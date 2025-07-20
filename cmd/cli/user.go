package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	v1 "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1/userv1connect"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/server"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user/store/sqlite"
)

var (
	apiEndpoint string
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
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
	RootCmd.AddCommand(userCmd)

	// Add API endpoint flag to the user command
	userCmd.PersistentFlags().StringVar(&apiEndpoint, "endpoint", "", "API endpoint URL (e.g., http://localhost:8088)")

	// Add all User RPC commands
	userCmd.AddCommand(listUsersCmd())
	userCmd.AddCommand(getUserCmd())
	userCmd.AddCommand(createUserCmd())
	userCmd.AddCommand(updateUserCmd())
	userCmd.AddCommand(deleteUserCmd())
}

// getClient returns either a local client or a remote client based on whether the API endpoint is provided
func getClient(ctx context.Context) (v1.UserServiceClient, error) {
	if apiEndpoint != "" {
		// Use Connect client with remote endpoint
		httpClient := http.DefaultClient
		return v1.NewUserServiceClient(
			httpClient,
			apiEndpoint,
		), nil
	} else {
		store, err := sqlite.NewStore(ctx, ":memory:")
		if err != nil {
			slog.DebugContext(ctx, "could not get sqlite user store", slog.Any("error", err))
			return nil, err
		}

		// Use local service with ServiceAdapter
		return server.NewUserConnectHandler(user.NewService(store)), nil
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
