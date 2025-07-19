package api

import (
	"context"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1"
)

func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// TODO: Implement actual logic
	return &pb.DeleteUserResponse{}, nil
}
