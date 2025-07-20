package server

import (
	"context"

	"connectrpc.com/connect"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
	"github.com/andrew-womeldorf/connect-boilerplate/pkg/api"
)

// ConnectHandler handles the over-the-wire connect requests, and sends them to
// the service, which handles in-memory objects.
type ConnectHandler struct {
	service *api.Service
}

// NewConnectHandler creates a new service adapter
func NewConnectHandler(service *api.Service) *ConnectHandler {
	return &ConnectHandler{
		service: service,
	}
}

// ListUsers implements the Connect interface
func (a *ConnectHandler) ListUsers(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
	resp, err := a.service.ListUsers(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetUser implements the Connect interface
func (a *ConnectHandler) GetUser(ctx context.Context, req *connect.Request[pb.GetUserRequest]) (*connect.Response[pb.GetUserResponse], error) {
	resp, err := a.service.GetUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// CreateUser implements the Connect interface
func (a *ConnectHandler) CreateUser(ctx context.Context, req *connect.Request[pb.CreateUserRequest]) (*connect.Response[pb.CreateUserResponse], error) {
	resp, err := a.service.CreateUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateUser implements the Connect interface
func (a *ConnectHandler) UpdateUser(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error) {
	resp, err := a.service.UpdateUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteUser implements the Connect interface
func (a *ConnectHandler) DeleteUser(ctx context.Context, req *connect.Request[pb.DeleteUserRequest]) (*connect.Response[pb.DeleteUserResponse], error) {
	resp, err := a.service.DeleteUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
