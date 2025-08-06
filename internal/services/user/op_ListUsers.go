package user

import (
	"context"
	"log/slog"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func (s *Service) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	slog.InfoContext(ctx, "listing users")

	users, err := s.store.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListUsersResponse{Users: users}, nil
}
