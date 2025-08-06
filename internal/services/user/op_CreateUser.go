package user

import (
	"context"
	"log/slog"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
	"github.com/google/uuid"
)

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	slog.InfoContext(ctx, "creating user", slog.String("name", req.Name), slog.String("email", req.Email))

	user := &pb.User{
		Id:    uuid.New().String(),
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{User: user}, nil
}
