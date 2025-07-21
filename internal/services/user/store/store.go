package store

import (
	"context"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

type Store interface {
	CreateUser(context.Context, *pb.User) error
	DeleteUser(context.Context, string) error
	GetUser(context.Context, string) (*pb.User, error)
	ListUsers(context.Context) ([]*pb.User, error)
	UpdateUser(context.Context, *pb.User) error
}
