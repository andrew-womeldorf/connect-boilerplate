package user

import (
	"context"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{User: user}, nil
}
