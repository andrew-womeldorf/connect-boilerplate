package server

import (
	"context"

	"connectrpc.com/connect"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1"
	"github.com/andrew-womeldorf/connect-boilerplate/pkg/api"
)

// ServiceAdapter adapts the API service to the Connect interface
type ServiceAdapter struct {
	service *api.Service
}

// NewServiceAdapter creates a new service adapter
func NewServiceAdapter(service *api.Service) *ServiceAdapter {
	return &ServiceAdapter{
		service: service,
	}
}

// ListUsers implements the Connect interface
func (a *ServiceAdapter) ListUsers(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
	return a.service.ListUsers(ctx, req)
}

// GetUser implements the Connect interface
func (a *ServiceAdapter) GetUser(ctx context.Context, req *connect.Request[pb.GetUserRequest]) (*connect.Response[pb.GetUserResponse], error) {
	return a.service.GetUser(ctx, req)
}

// CreateUser implements the Connect interface
func (a *ServiceAdapter) CreateUser(ctx context.Context, req *connect.Request[pb.CreateUserRequest]) (*connect.Response[pb.CreateUserResponse], error) {
	return a.service.CreateUser(ctx, req)
}

// UpdateUser implements the Connect interface
func (a *ServiceAdapter) UpdateUser(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error) {
	return a.service.UpdateUser(ctx, req)
}

// DeleteUser implements the Connect interface
func (a *ServiceAdapter) DeleteUser(ctx context.Context, req *connect.Request[pb.DeleteUserRequest]) (*connect.Response[pb.DeleteUserResponse], error) {
	return a.service.DeleteUser(ctx, req)
}
