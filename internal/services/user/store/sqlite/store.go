package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	_ "modernc.org/sqlite"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user/store/sqlite/gen"
)

var (
	ErrCouldNotGetUser    = errors.New("could not get user")
	ErrCouldNotCreateUser = errors.New("could not create user")
	ErrCouldNotDeleteUser = errors.New("could not delete user")
	ErrCouldNotListUsers  = errors.New("could not list users")
	ErrCouldNotUpdateUser = errors.New("could not update user")
)

//go:embed schema.sql
var Schema string

type Store struct {
	q *gen.Queries
}

func NewStore(ctx context.Context, sqliteFile string) (*Store, error) {
	db, err := sql.Open("sqlite", sqliteFile)
	if err != nil {
		return nil, err
	}

	return &Store{
		q: gen.New(db),
	}, nil
}

func (s *Store) CreateUser(ctx context.Context, user *pb.User) error {
	if _, err := s.q.CreateUser(ctx, gen.CreateUserParams{
		ID:        user.GetId(),
		Name:      user.GetName(),
		Email:     user.GetEmail(),
		CreatedAt: user.GetCreatedAt().AsTime().Format(time.RFC3339),
		UpdatedAt: user.GetUpdatedAt().AsTime().Format(time.RFC3339),
	}); err != nil {
		slog.ErrorContext(ctx, ErrCouldNotCreateUser.Error(),
			slog.Any("error", err),
			slog.String("user id", user.GetId()),
		)
		return ErrCouldNotCreateUser
	}

	return nil
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	_, err := s.q.DeleteUser(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotDeleteUser.Error(),
			slog.Any("error", err),
			slog.String("user id", id),
		)
		return ErrCouldNotDeleteUser
	}

	return nil
}

func (s *Store) GetUser(ctx context.Context, id string) (*pb.User, error) {
	db, err := s.q.GetUser(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotGetUser.Error(),
			slog.Any("error", err),
			slog.String("user id", id),
		)
		return nil, ErrCouldNotGetUser
	}

	user, err := convertUser(ctx, db)
	if err != nil {
		return nil, ErrCouldNotGetUser
	}

	return user, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]*pb.User, error) {
	db, err := s.q.ListUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotListUsers.Error(),
			slog.Any("error", err),
		)
		return nil, ErrCouldNotGetUser
	}

	users := []*pb.User{}
	for _, u := range db {
		pbu, err := convertUser(ctx, u)
		if err != nil {
			return nil, ErrCouldNotListUsers
		}
		users = append(users, pbu)
	}

	return users, nil
}

func (s *Store) UpdateUser(ctx context.Context, user *pb.User) error {
	_, err := s.q.UpdateUser(ctx, gen.UpdateUserParams{
		ID:        user.GetId(),
		Name:      user.GetName(),
		Email:     user.GetEmail(),
		UpdatedAt: user.GetUpdatedAt().AsTime().Format(time.RFC3339),
	})

	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotUpdateUser.Error(),
			slog.Any("error", err),
			slog.String("user id", user.GetId()),
		)
		return ErrCouldNotUpdateUser
	}

	return nil
}

func convertUser(ctx context.Context, db gen.User) (*pb.User, error) {
	user := &pb.User{
		Id:    db.ID,
		Name:  db.Name,
		Email: db.Email,
	}

	t, err := time.Parse(time.RFC3339, db.CreatedAt)
	if err != nil {
		msg := "could not parse created at timestamp"
		slog.ErrorContext(ctx, msg,
			slog.Any("error", err),
			slog.String("user id", db.ID),
			slog.String("created at", db.CreatedAt),
		)
		return nil, fmt.Errorf("%s: %w", msg, err)
	}
	user.CreatedAt = timestamppb.New(t)

	t, err = time.Parse(time.RFC3339, db.UpdatedAt)
	if err != nil {
		msg := "could not parse updated at timestamp"
		slog.ErrorContext(ctx, msg,
			slog.Any("error", err),
			slog.String("user id", db.ID),
			slog.String("updated at", db.UpdatedAt),
		)
		return nil, fmt.Errorf("%s: %w", msg, err)
	}
	user.UpdatedAt = timestamppb.New(t)

	return user, nil
}
