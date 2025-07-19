package api

import (
	"context"
	"log/slog"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1"
)

// Service handles the business logic
type Service struct {
	// Add dependencies here (database, other services, etc.)
}

func NewService() *Service {
	return &Service{}
}

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

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// TODO: Implement actual logic
	user := &pb.User{
		Id:    req.Id,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	return &pb.GetUserResponse{User: user}, nil
}

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

func (s *Service) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO: Implement actual logic
	user := &pb.User{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
	}

	return &pb.UpdateUserResponse{User: user}, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// TODO: Implement actual logic
	return &pb.DeleteUserResponse{}, nil
}
