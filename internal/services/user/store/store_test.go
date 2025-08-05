package store

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	tc "github.com/testcontainers/testcontainers-go/modules/dynamodb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
	ddbstore "github.com/andrew-womeldorf/connect-boilerplate/internal/services/user/store/dynamodb"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user/store/sqlite"
)

type storeTestSuite struct {
	name  string
	setup func(t *testing.T) (Store, func())
}

type ddbResolver struct {
	port string
}

func (r *ddbResolver) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return smithyendpoints.Endpoint{URI: url.URL{Host: r.port, Scheme: "http"}}, nil
}

var (
	sharedDynamoDBContainer *tc.DynamoDBContainer
	sharedDynamoDBClient    *dynamodb.Client
	sharedDynamoDBTableName = "users"
	containerSetupOnce      sync.Once
)

func setupSharedDynamoDBContainer() error {
	var err error
	containerSetupOnce.Do(func() {
		ctx := context.Background()

		sharedDynamoDBContainer, err = tc.Run(ctx, "amazon/dynamodb-local:latest", tc.WithSharedDB())
		if err != nil {
			err = fmt.Errorf("could not start dynamodb container: %w", err)
			return
		}

		port, portErr := sharedDynamoDBContainer.ConnectionString(ctx)
		if portErr != nil {
			err = fmt.Errorf("could not get connection string from dynamodb container: %w", portErr)
			return
		}

		cfg, cfgErr := config.LoadDefaultConfig(ctx,
			config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{AccessKeyID: "dummy", SecretAccessKey: "dummy"},
			}),
		)
		if cfgErr != nil {
			err = fmt.Errorf("failed to create aws config: %w", cfgErr)
			return
		}

		sharedDynamoDBClient = dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolverV2(&ddbResolver{port: port}))

		_, tableErr := sharedDynamoDBClient.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(sharedDynamoDBTableName),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("PK"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("SK"),
					KeyType:       types.KeyTypeRange,
				},
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI1PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI1SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("GSI1"),
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("GSI1PK"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("GSI1SK"),
							KeyType:       types.KeyTypeRange,
						},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
				},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if tableErr != nil {
			err = fmt.Errorf("failed to create dynamodb table: %w", tableErr)
			return
		}
	})
	return err
}

func cleanupDynamoDBTable(ctx context.Context) error {
	if sharedDynamoDBClient == nil {
		return nil
	}

	scanOutput, err := sharedDynamoDBClient.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(sharedDynamoDBTableName),
	})
	if err != nil {
		return fmt.Errorf("failed to scan table for cleanup: %w", err)
	}

	for _, item := range scanOutput.Items {
		_, err := sharedDynamoDBClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(sharedDynamoDBTableName),
			Key: map[string]types.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete item during cleanup: %w", err)
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	code := m.Run()

	// Cleanup shared container after all tests
	if sharedDynamoDBContainer != nil {
		ctx := context.Background()
		if err := sharedDynamoDBContainer.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate shared dynamodb container: %v\n", err)
		}
	}

	os.Exit(code)
}

func TestStore(t *testing.T) {
	ctx := context.Background()

	testSuites := []storeTestSuite{
		{
			name: "SQLite",
			setup: func(t *testing.T) (Store, func()) {
				tmpFile, err := os.CreateTemp("", "test_*.db")
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
				if err = tmpFile.Close(); err != nil {
					t.Fatalf("failed to close temp file: %v", err)
				}

				db, err := sql.Open("sqlite", tmpFile.Name())
				if err != nil {
					t.Fatalf("failed to open database: %v", err)
				}

				if _, err := db.Exec(sqlite.Schema); err != nil {
					t.Fatalf("failed to create schema: %v", err)
				}
				if err = db.Close(); err != nil {
					t.Fatalf("failed to close sqlite file: %v", err)
				}

				store, err := sqlite.NewStore(ctx, tmpFile.Name())
				if err != nil {
					t.Fatalf("failed to create store: %v", err)
				}

				cleanup := func() {
					if err := os.Remove(tmpFile.Name()); err != nil {
						t.Fatal("failed to cleanup sqlite file")
					}
				}

				return store, cleanup
			},
		},
		{
			name: "DynamoDB",
			setup: func(t *testing.T) (Store, func()) {
				// Setup shared container if not already done
				if err := setupSharedDynamoDBContainer(); err != nil {
					t.Fatalf("failed to setup shared dynamodb container: %v", err)
				}

				// Clean the table before each test
				if err := cleanupDynamoDBTable(ctx); err != nil {
					t.Fatalf("failed to cleanup dynamodb table: %v", err)
				}

				// create the store using the shared client and table
				store, err := ddbstore.NewStore(
					ctx, ddbstore.WithClient(sharedDynamoDBClient), ddbstore.WithTable(sharedDynamoDBTableName),
				)
				if err != nil {
					t.Fatalf("failed to create store: %v", err)
				}

				cleanup := func() {
					// Clean the table after each test
					if err := cleanupDynamoDBTable(ctx); err != nil {
						t.Logf("failed to cleanup dynamodb table: %v", err)
					}
				}

				return store, cleanup
			},
		},
	}

	for _, suite := range testSuites {
		t.Run(suite.name, func(t *testing.T) {
			runStoreTests(ctx, t, suite.setup)
		})
	}
}

func runStoreTests(ctx context.Context, t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("CreateUser", func(t *testing.T) {
		testCreateUser(ctx, t, setup)
	})
	t.Run("GetUser", func(t *testing.T) {
		testGetUser(ctx, t, setup)
	})
	t.Run("UpdateUser", func(t *testing.T) {
		testUpdateUser(ctx, t, setup)
	})
	t.Run("DeleteUser", func(t *testing.T) {
		testDeleteUser(ctx, t, setup)
	})
	t.Run("ListUsers", func(t *testing.T) {
		testListUsers(ctx, t, setup)
	})
}

func testCreateUser(ctx context.Context, t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("success", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("1", "John Doe", "john@example.com")
		err := store.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		retrieved, err := store.GetUser(ctx, "1")
		if err != nil {
			t.Fatalf("failed to retrieve created user: %v", err)
		}

		if retrieved.GetId() != user.GetId() {
			t.Errorf("expected id %s, got %s", user.GetId(), retrieved.GetId())
		}
		if retrieved.GetName() != user.GetName() {
			t.Errorf("expected name %s, got %s", user.GetName(), retrieved.GetName())
		}
		if retrieved.GetEmail() != user.GetEmail() {
			t.Errorf("expected email %s, got %s", user.GetEmail(), retrieved.GetEmail())
		}
	})

	t.Run("duplicate_id", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user1 := createTestUser("1", "John Doe", "john@example.com")
		user2 := createTestUser("1", "Jane Doe", "jane@example.com")

		err := store.CreateUser(ctx, user1)
		if err != nil {
			t.Fatalf("first create should succeed: %v", err)
		}

		err = store.CreateUser(ctx, user2)
		if err == nil {
			t.Fatal("expected error for duplicate ID, got nil")
		}
	})
}

func testGetUser(ctx context.Context, t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("1", "John Doe", "john@example.com")
		err := store.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		retrieved, err := store.GetUser(ctx, "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if retrieved.GetId() != user.GetId() {
			t.Errorf("expected id %s, got %s", user.GetId(), retrieved.GetId())
		}
		if retrieved.GetName() != user.GetName() {
			t.Errorf("expected name %s, got %s", user.GetName(), retrieved.GetName())
		}
		if retrieved.GetEmail() != user.GetEmail() {
			t.Errorf("expected email %s, got %s", user.GetEmail(), retrieved.GetEmail())
		}
	})

	t.Run("non_existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		_, err := store.GetUser(ctx, "non-existent")
		if err == nil {
			t.Fatal("expected error for non-existent user, got nil")
		}
	})
}

func testUpdateUser(ctx context.Context, t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("1", "John Doe", "john@example.com")
		err := store.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		updatedUser := createTestUser("1", "John Smith", "johnsmith@example.com")
		err = store.UpdateUser(ctx, updatedUser)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		retrieved, err := store.GetUser(ctx, "1")
		if err != nil {
			t.Fatalf("failed to retrieve updated user: %v", err)
		}

		if retrieved.GetName() != "John Smith" {
			t.Errorf("expected updated name 'John Smith', got %s", retrieved.GetName())
		}
		if retrieved.GetEmail() != "johnsmith@example.com" {
			t.Errorf("expected updated email 'johnsmith@example.com', got %s", retrieved.GetEmail())
		}
	})

	t.Run("non_existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("non-existent", "John Doe", "john@example.com")
		err := store.UpdateUser(ctx, user)
		if err == nil {
			t.Fatal("expected error for non-existent user, got nil")
		}
	})
}

func testDeleteUser(ctx context.Context, t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("1", "John Doe", "john@example.com")
		err := store.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		err = store.DeleteUser(ctx, "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = store.GetUser(ctx, "1")
		if err == nil {
			t.Fatal("user should not exist after deletion")
		}
	})

	t.Run("non_existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		err := store.DeleteUser(ctx, "non-existent")
		if err == nil {
			t.Fatalf("deleting non-existent user should fail, got %v", err)
		}
	})
}

func testListUsers(ctx context.Context, t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("empty_list", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		users, err := store.ListUsers(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(users) != 0 {
			t.Errorf("expected empty list, got %d users", len(users))
		}
	})

	t.Run("multiple_users", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user1 := createTestUser("1", "John Doe", "john@example.com")
		user2 := createTestUser("2", "Jane Doe", "jane@example.com")
		user3 := createTestUser("3", "Bob Smith", "bob@example.com")

		for _, user := range []*pb.User{user1, user2, user3} {
			err := store.CreateUser(ctx, user)
			if err != nil {
				t.Fatalf("failed to create user %s: %v", user.GetId(), err)
			}
		}

		users, err := store.ListUsers(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(users) != 3 {
			t.Errorf("expected 3 users, got %d", len(users))
		}

		userMap := make(map[string]*pb.User)
		for _, user := range users {
			userMap[user.GetId()] = user
		}

		expectedUsers := []*pb.User{user1, user2, user3}
		for _, expected := range expectedUsers {
			actual, exists := userMap[expected.GetId()]
			if !exists {
				t.Errorf("user with ID %s not found in list", expected.GetId())
				continue
			}

			if actual.GetName() != expected.GetName() {
				t.Errorf("expected name %s, got %s for user %s", expected.GetName(), actual.GetName(), expected.GetId())
			}
			if actual.GetEmail() != expected.GetEmail() {
				t.Errorf("expected email %s, got %s for user %s", expected.GetEmail(), actual.GetEmail(), expected.GetId())
			}
		}
	})
}

func createTestUser(id, name, email string) *pb.User {
	now := time.Now()
	return &pb.User{
		Id:        id,
		Name:      name,
		Email:     email,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}
}
