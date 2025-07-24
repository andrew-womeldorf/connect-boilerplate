package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
)

const defaultTableName = "users"

var (
	ErrCouldNotGetUser    = errors.New("could not get user")
	ErrCouldNotCreateUser = errors.New("could not create user")
	ErrCouldNotDeleteUser = errors.New("could not delete user")
	ErrCouldNotListUsers  = errors.New("could not list users")
	ErrCouldNotUpdateUser = errors.New("could not update user")
)

type Store struct {
	client *ddb.Client
	table  string
}

type Option func(*Store)

func WithTable(name string) Option {
	return func(s *Store) {
		s.table = name
	}
}

func WithClient(client *ddb.Client) Option {
	return func(s *Store) {
		s.client = client
	}
}

func NewStore(ctx context.Context, opts ...Option) (*Store, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "could not load default aws config", slog.Any("error", err))
		return nil, err
	}

	s := &Store{
		client: ddb.NewFromConfig(cfg),
		table:  defaultTableName,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

type User struct {
	Id        string    `dynamodbav:"id"`
	Name      string    `dynamodbav:"name"`
	Email     string    `dynamodbav:"email"`
	CreatedAt time.Time `dynamodbav:"createdAt"`
	UpdatedAt time.Time `dynamodbav:"updatedAt"`
}

type UserItem struct {
	PK     string `dynamodbav:"PK"`
	SK     string `dynamodbav:"SK"`
	GSI1PK string `dynamodbav:"GSI1PK"`
	GSI1SK string `dynamodbav:"GSI1SK"`
	User   User   `dynamodbav:"user"`
}

func (item *UserItem) SetKeys() {
	item.PK = fmt.Sprintf("USER#%s", item.User.Id)
	item.SK = fmt.Sprintf("USER#%s", item.User.Id)
	item.GSI1PK = "USERS"
	item.GSI1SK = item.User.Id
}

func (s *Store) CreateUser(ctx context.Context, user *pb.User) error {
	item := UserItem{
		User: User{
			Id:        user.GetId(),
			Name:      user.GetName(),
			Email:     user.GetEmail(),
			CreatedAt: user.GetCreatedAt().AsTime(),
			UpdatedAt: user.GetUpdatedAt().AsTime(),
		},
	}
	item.SetKeys()

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotCreateUser.Error(),
			slog.Any("error", err),
			slog.String("user id", user.GetId()),
		)
		return ErrCouldNotCreateUser
	}

	_, err = s.client.PutItem(ctx, &ddb.PutItemInput{
		TableName:           &s.table,
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})

	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotCreateUser.Error(),
			slog.Any("error", err),
			slog.String("user id", user.GetId()),
		)
		return ErrCouldNotCreateUser
	}

	return nil
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", id)},
		"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", id)},
	}

	_, err := s.client.DeleteItem(ctx, &ddb.DeleteItemInput{
		TableName:           &s.table,
		Key:                 key,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})

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
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", id)},
		"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", id)},
	}

	resp, err := s.client.GetItem(ctx, &ddb.GetItemInput{
		TableName: &s.table,
		Key:       key,
	})

	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotGetUser.Error(),
			slog.Any("error", err),
			slog.String("user id", id),
		)
		return nil, ErrCouldNotGetUser
	}

	if resp.Item == nil {
		return nil, ErrCouldNotGetUser
	}

	var item UserItem
	if err := attributevalue.UnmarshalMap(resp.Item, &item); err != nil {
		slog.ErrorContext(ctx, ErrCouldNotGetUser.Error(),
			slog.Any("error", err),
			slog.String("user id", id),
		)
		return nil, ErrCouldNotGetUser
	}

	return convertUserItem(item), nil
}

func (s *Store) ListUsers(ctx context.Context) ([]*pb.User, error) {
	resp, err := s.client.Query(ctx, &ddb.QueryInput{
		TableName:              &s.table,
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USERS"},
		},
	})

	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotListUsers.Error(),
			slog.Any("error", err),
		)
		return nil, ErrCouldNotListUsers
	}

	users := make([]*pb.User, 0, len(resp.Items))
	for _, item := range resp.Items {
		var userItem UserItem
		if err := attributevalue.UnmarshalMap(item, &userItem); err != nil {
			slog.ErrorContext(ctx, ErrCouldNotListUsers.Error(),
				slog.Any("error", err),
			)
			continue
		}
		users = append(users, convertUserItem(userItem))
	}

	return users, nil
}

func (s *Store) UpdateUser(ctx context.Context, user *pb.User) error {
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", user.GetId())},
		"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", user.GetId())},
	}

	userData := User{
		Id:        user.GetId(),
		Name:      user.GetName(),
		Email:     user.GetEmail(),
		CreatedAt: user.GetCreatedAt().AsTime(),
		UpdatedAt: user.GetUpdatedAt().AsTime(),
	}

	userAv, err := attributevalue.Marshal(userData)
	if err != nil {
		slog.ErrorContext(ctx, ErrCouldNotUpdateUser.Error(),
			slog.Any("error", err),
			slog.String("user id", user.GetId()),
		)
		return ErrCouldNotUpdateUser
	}

	_, err = s.client.UpdateItem(ctx, &ddb.UpdateItemInput{
		TableName:        &s.table,
		Key:              key,
		UpdateExpression: aws.String("SET #user = :user"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user": userAv,
		},
		ExpressionAttributeNames: map[string]string{
			"#user": "user",
		},
		ConditionExpression: aws.String("attribute_exists(PK)"),
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

func convertUserItem(item UserItem) *pb.User {
	return &pb.User{
		Id:        item.User.Id,
		Name:      item.User.Name,
		Email:     item.User.Email,
		CreatedAt: timestamppb.New(item.User.CreatedAt),
		UpdatedAt: timestamppb.New(item.User.UpdatedAt),
	}
}
