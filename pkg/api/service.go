package api

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1"
)

type Service struct {
	// Add dependencies here (database, other services, etc.)
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ListUsers(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
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

	resp := connect.NewResponse(&pb.ListUsersResponse{
		Users: users,
	})

	return resp, nil
}

func (s *Service) GetUser(ctx context.Context, req *connect.Request[pb.GetUserRequest]) (*connect.Response[pb.GetUserResponse], error) {
	// TODO: Implement actual logic
	user := &pb.User{
		Id:    req.Msg.Id,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	resp := connect.NewResponse(&pb.GetUserResponse{
		User: user,
	})

	return resp, nil
}

func (s *Service) CreateUser(ctx context.Context, req *connect.Request[pb.CreateUserRequest]) (*connect.Response[pb.CreateUserResponse], error) {
	slog.InfoContext(ctx, "creating user", slog.String("name", req.Msg.Name), slog.String("email", req.Msg.Email))

	// TODO: Implement actual logic
	user := &pb.User{
		Id:    "new-user-id",
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
	}

	resp := connect.NewResponse(&pb.CreateUserResponse{
		User: user,
	})

	return resp, nil
}

func (s *Service) UpdateUser(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error) {
	// TODO: Implement actual logic
	user := &pb.User{
		Id:    req.Msg.Id,
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
	}

	resp := connect.NewResponse(&pb.UpdateUserResponse{
		User: user,
	})

	return resp, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *connect.Request[pb.DeleteUserRequest]) (*connect.Response[pb.DeleteUserResponse], error) {
	// TODO: Implement actual logic
	resp := connect.NewResponse(&pb.DeleteUserResponse{})

	return resp, nil
}
