package user

import (
	"context"
	"log/slog"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	slog.InfoContext(ctx, "creating user", slog.String("name", req.Name), slog.String("email", req.Email))

	// TODO: Implement actual logic
	user := &pb.User{
		Id:    "new-user-id",
		Name:  req.Name,
		Email: req.Email,
	}

	return &pb.CreateUserResponse{User: user}, nil
}
