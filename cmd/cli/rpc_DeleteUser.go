package main

import (
	"context"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1"
)

func deleteUserCmd() *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "delete-user",
		Short: "Delete a user by ID",
		Long:  `Delete a user by their ID.`,
		Run: func(cmd *cobra.Command, args []string) {
			runDeleteUser(userID)
		},
	}

	cmd.Flags().StringVar(&userID, "id", "", "User ID to delete (required)")
	if err := cmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}

	return cmd
}

func runDeleteUser(userID string) {
	ctx := context.Background()

	// Get client based on endpoint flag
	client, err := getClient(ctx)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}

	// Create request
	req := &pb.DeleteUserRequest{
		Id: userID,
	}

	// Call the service
	slog.Debug("Deleting user", "id", userID)
	resp, err := client.DeleteUser(ctx, connect.NewRequest(req))
	if err != nil {
		slog.Error("Failed to delete user", "error", err)
		os.Exit(1)
	}
	slog.Debug("Successfully deleted user")

	printJSON(resp.Msg)
}
