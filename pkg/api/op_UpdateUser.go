package api

import (
	"context"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1"
)

func (s *Service) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO: Implement actual logic
	user := &pb.User{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
	}

	return &pb.UpdateUserResponse{User: user}, nil
}
