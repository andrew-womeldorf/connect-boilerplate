package user

import (
	"context"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.store.DeleteUser(ctx, req.Id); err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{}, nil
}
