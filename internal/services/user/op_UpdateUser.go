package user

import (
	"context"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

func (s *Service) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := &pb.User{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.store.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{User: user}, nil
}
