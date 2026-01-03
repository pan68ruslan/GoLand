package users

import (
	"testing"

	store "github.com/pan68ruslan/GoLand/lesson_05/documentstore"
)

var testStore = "TestDB"
var testService = "Users"
var testUser = "User0"
var testId = "u0"

func newTestService(t *testing.T) *Service {
	s := store.NewStore("")
	if s != nil && (s.Name != "Noname_Store" || s.Collections == nil) {
		t.Fatal("failed to create the empty service")
	}
	s = store.NewStore(testStore)
	service, err := NewService(testService, s)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	return service
}

func TestCreateUser(t *testing.T) {
	service := newTestService(t)
	created, err := service.CreateUser(User{ID: testId, Name: testUser})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if created.ID != testId || created.Name != testUser {
		t.Fatalf("Unexpected user data: %+v", created)
	}
	_, err = service.CreateUser(User{ID: testId, Name: testUser})
	if err == nil {
		t.Fatalf("Expected error for duplicate user, got nil")
	}
}

func TestGetUser(t *testing.T) {
	service := newTestService(t)
	_, _ = service.CreateUser(User{ID: testId, Name: testUser})
	got, err := service.GetUser(testId)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if got.ID != testId || got.Name != testUser {
		t.Fatalf("Expected Alice, got %+v", got)
	}
	_, err = service.GetUser("unknown")
	if err == nil {
		t.Errorf("Expected error for unknown user, got nil")
	}
}

func TestDeleteUser(t *testing.T) {
	service := newTestService(t)
	u := User{ID: testId, Name: testUser}
	_, _ = service.CreateUser(u)
	if err := service.DeleteUser(testId); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	if err := service.DeleteUser(testId); err == nil {
		t.Errorf("expected error when deleting non-existing user, got nil")
	}
}

func TestListUsers(t *testing.T) {
	service := newTestService(t)
	users, err := service.ListUsers()
	if err != nil {
		t.Errorf("expected no error for empty list, got %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
	_, _ = service.CreateUser(User{ID: "u1", Name: "User1"})
	_, _ = service.CreateUser(User{ID: "u2", Name: "User2"})
	users, err = service.ListUsers()
	if err != nil {
		t.Errorf("ListUsers returned error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}
