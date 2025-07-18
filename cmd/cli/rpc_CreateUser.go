package main

import (
	"context"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1"
)

func createUserCmd() *cobra.Command {
	var userName string
	var userEmail string

	cmd := &cobra.Command{
		Use:   "create-user",
		Short: "Create a new user",
		Long:  `Create a new user with the given name and email.`,
		Run: func(cmd *cobra.Command, args []string) {
			runCreateUser(userName, userEmail)
		},
	}

	cmd.Flags().StringVar(&userName, "name", "", "User name (required)")
	cmd.Flags().StringVar(&userEmail, "email", "", "User email (required)")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("email"); err != nil {
		panic(err)
	}

	return cmd
}

func runCreateUser(userName, userEmail string) {
	ctx := context.Background()

	// Get client based on endpoint flag
	client, err := getClient(ctx)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}

	// Create request
	req := &pb.CreateUserRequest{
		Name:  userName,
		Email: userEmail,
	}

	// Call the service
	slog.Debug("Creating user", "name", userName, "email", userEmail)
	resp, err := client.CreateUser(ctx, connect.NewRequest(req))
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		os.Exit(1)
	}
	slog.Debug("Successfully created user")

	printJSON(resp.Msg)
}
