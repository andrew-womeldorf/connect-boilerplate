package store

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user/store/sqlite"
)

type storeTestSuite struct {
	name  string
	setup func(t *testing.T) (Store, func())
}

func TestStore(t *testing.T) {
	testSuites := []storeTestSuite{
		{
			name: "SQLite",
			setup: func(t *testing.T) (Store, func()) {
				tmpFile, err := os.CreateTemp("", "test_*.db")
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
				tmpFile.Close()

				db, err := sql.Open("sqlite", tmpFile.Name())
				if err != nil {
					t.Fatalf("failed to open database: %v", err)
				}

				if _, err := db.Exec(sqlite.Schema); err != nil {
					t.Fatalf("failed to create schema: %v", err)
				}
				db.Close()

				store, err := sqlite.NewStore(context.Background(), tmpFile.Name())
				if err != nil {
					t.Fatalf("failed to create store: %v", err)
				}

				cleanup := func() {
					os.Remove(tmpFile.Name())
				}

				return store, cleanup
			},
		},
	}

	for _, suite := range testSuites {
		t.Run(suite.name, func(t *testing.T) {
			runStoreTests(t, suite.setup)
		})
	}
}

func runStoreTests(t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("CreateUser", func(t *testing.T) {
		testCreateUser(t, setup)
	})
	t.Run("GetUser", func(t *testing.T) {
		testGetUser(t, setup)
	})
	t.Run("UpdateUser", func(t *testing.T) {
		testUpdateUser(t, setup)
	})
	t.Run("DeleteUser", func(t *testing.T) {
		testDeleteUser(t, setup)
	})
	t.Run("ListUsers", func(t *testing.T) {
		testListUsers(t, setup)
	})
}

func testCreateUser(t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("success", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("1", "John Doe", "john@example.com")
		err := store.CreateUser(context.Background(), user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		retrieved, err := store.GetUser(context.Background(), "1")
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

		err := store.CreateUser(context.Background(), user1)
		if err != nil {
			t.Fatalf("first create should succeed: %v", err)
		}

		err = store.CreateUser(context.Background(), user2)
		if err == nil {
			t.Fatal("expected error for duplicate ID, got nil")
		}
	})
}

func testGetUser(t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("1", "John Doe", "john@example.com")
		err := store.CreateUser(context.Background(), user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		retrieved, err := store.GetUser(context.Background(), "1")
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

		_, err := store.GetUser(context.Background(), "non-existent")
		if err == nil {
			t.Fatal("expected error for non-existent user, got nil")
		}
	})
}

func testUpdateUser(t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("1", "John Doe", "john@example.com")
		err := store.CreateUser(context.Background(), user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		updatedUser := createTestUser("1", "John Smith", "johnsmith@example.com")
		err = store.UpdateUser(context.Background(), updatedUser)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		retrieved, err := store.GetUser(context.Background(), "1")
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
		err := store.UpdateUser(context.Background(), user)
		if err == nil {
			t.Fatal("expected error for non-existent user, got nil")
		}
	})
}

func testDeleteUser(t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		user := createTestUser("1", "John Doe", "john@example.com")
		err := store.CreateUser(context.Background(), user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		err = store.DeleteUser(context.Background(), "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = store.GetUser(context.Background(), "1")
		if err == nil {
			t.Fatal("user should not exist after deletion")
		}
	})

	t.Run("non_existing_user", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		err := store.DeleteUser(context.Background(), "non-existent")
		if err == nil {
			t.Fatalf("deleting non-existent user should fail, got %v", err)
		}
	})
}

func testListUsers(t *testing.T, setup func(t *testing.T) (Store, func())) {
	t.Run("empty_list", func(t *testing.T) {
		store, cleanup := setup(t)
		defer cleanup()

		users, err := store.ListUsers(context.Background())
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
			err := store.CreateUser(context.Background(), user)
			if err != nil {
				t.Fatalf("failed to create user %s: %v", user.GetId(), err)
			}
		}

		users, err := store.ListUsers(context.Background())
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
