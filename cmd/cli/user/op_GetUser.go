package user

import (
	"context"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func getUserCmd() *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "get-user",
		Short: "Get a user by ID",
		Long:  `Get a user by their ID.`,
		Run: func(cmd *cobra.Command, args []string) {
			runGetUser(userID)
		},
	}

	cmd.Flags().StringVar(&userID, "id", "", "User ID to retrieve (required)")
	if err := cmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}

	return cmd
}

func runGetUser(userID string) {
	ctx := context.Background()

	// Get client based on endpoint flag
	client, err := getClient(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create client", "error", err)
		os.Exit(1)
	}

	// Create request
	req := &pb.GetUserRequest{
		Id: userID,
	}

	// Call the service
	slog.DebugContext(ctx, "Getting user", "id", userID)
	resp, err := client.GetUser(ctx, connect.NewRequest(req))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get user", "error", err)
		os.Exit(1)
	}
	slog.DebugContext(ctx, "Successfully got user")

	printJSON(resp.Msg)
}
