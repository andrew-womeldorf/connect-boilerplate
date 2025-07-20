package main

import (
	"context"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func updateUserCmd() *cobra.Command {
	var userID string
	var userName string
	var userEmail string

	cmd := &cobra.Command{
		Use:   "update-user",
		Short: "Update an existing user",
		Long:  `Update an existing user with the given ID, name, and email.`,
		Run: func(cmd *cobra.Command, args []string) {
			runUpdateUser(userID, userName, userEmail)
		},
	}

	cmd.Flags().StringVar(&userID, "id", "", "User ID to update (required)")
	cmd.Flags().StringVar(&userName, "name", "", "User name (required)")
	cmd.Flags().StringVar(&userEmail, "email", "", "User email (required)")
	if err := cmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("email"); err != nil {
		panic(err)
	}

	return cmd
}

func runUpdateUser(userID, userName, userEmail string) {
	ctx := context.Background()

	// Get client based on endpoint flag
	client, err := getClient(ctx)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}

	// Create request
	req := &pb.UpdateUserRequest{
		Id:    userID,
		Name:  userName,
		Email: userEmail,
	}

	// Call the service
	slog.Debug("Updating user", "id", userID, "name", userName, "email", userEmail)
	resp, err := client.UpdateUser(ctx, connect.NewRequest(req))
	if err != nil {
		slog.Error("Failed to update user", "error", err)
		os.Exit(1)
	}
	slog.Debug("Successfully updated user")

	printJSON(resp.Msg)
}
