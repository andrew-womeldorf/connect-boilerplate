package user

import (
	"context"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// TODO: Implement actual logic
	user := &pb.User{
		Id:    req.Id,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	return &pb.GetUserResponse{User: user}, nil
}
