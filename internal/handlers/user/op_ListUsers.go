package user

import (
	"context"
	"log/slog"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func (s *Service) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	slog.InfoContext(ctx, "listing users")

	// TODO: Implement actual logic
	users := []*pb.User{
		{
			Id:    "1",
			Name:  "John Doe",
			Email: "john@example.com",
		},
		{
			Id:    "2",
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
	}

	return &pb.ListUsersResponse{Users: users}, nil
}
