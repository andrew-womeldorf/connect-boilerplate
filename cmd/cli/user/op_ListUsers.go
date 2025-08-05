package user

import (
	"context"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func listUsersCmd() *cobra.Command {
	var pageSize int32
	var pageToken string

	cmd := &cobra.Command{
		Use:   "list-users",
		Short: "List users",
		Long: `List users with optional pagination.
This command allows listing users with page size and token parameters.`,
		Run: func(cmd *cobra.Command, args []string) {
			runListUsers(pageSize, pageToken)
		},
	}

	cmd.Flags().Int32Var(&pageSize, "page-size", 10, "Number of users to return per page")
	cmd.Flags().StringVar(&pageToken, "page-token", "", "Page token for pagination")

	return cmd
}

func runListUsers(pageSize int32, pageToken string) {
	ctx := context.Background()

	// Get client based on endpoint flag
	client, err := getClient(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create client", "error", err)
		os.Exit(1)
	}

	// Create request
	req := &pb.ListUsersRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	}

	// Call the service
	slog.DebugContext(ctx, "Listing users...")
	resp, err := client.ListUsers(ctx, connect.NewRequest(req))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list users", "error", err)
		os.Exit(1)
	}
	slog.DebugContext(ctx, "Successfully listed users", "count", len(resp.Msg.Users))

	// Create a response object with users
	result := struct {
		Users         []*pb.User `json:"users"`
		NextPageToken string     `json:"next_page_token,omitempty"`
	}{
		Users:         resp.Msg.Users,
		NextPageToken: resp.Msg.NextPageToken,
	}

	printJSON(result)
}
