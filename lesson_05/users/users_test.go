package users

import (
	"testing"

	store "github.com/pan68ruslan/GoLand/lesson_05/documentstore"
)

func newTestService(t *testing.T) *Service {
	s := store.NewStore("TestDB")
	service, err := NewService("users", s)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	return service
}

func TestCreateUser(t *testing.T) {
	service := newTestService(t)
	u := User{ID: "u1", Name: "Ruslan"}
	created, err := service.CreateUser(u)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if created.ID != u.ID || created.Name != u.Name {
		t.Errorf("expected %+v, got %+v", u, created)
	}
}

func TestGetUser(t *testing.T) {
	service := newTestService(t)
	u := User{ID: "u2", Name: "Anna"}
	_, _ = service.CreateUser(u)
	got, err := service.GetUser("u2")
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if got.ID != u.ID || got.Name != u.Name {
		t.Errorf("expected %+v, got %+v", u, got)
	}
	_, err = service.GetUser("unknown")
	if err == nil {
		t.Errorf("expected error for unknown user, got nil")
	}
}

func TestDeleteUser(t *testing.T) {
	service := newTestService(t)
	u := User{ID: "u3", Name: "Ivan"}
	_, _ = service.CreateUser(u)
	if err := service.DeleteUser("u3"); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	if err := service.DeleteUser("u3"); err == nil {
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
	_, _ = service.CreateUser(User{ID: "u4", Name: "Petro"})
	_, _ = service.CreateUser(User{ID: "u5", Name: "Oksana"})
	users, err = service.ListUsers()
	if err != nil {
		t.Errorf("ListUsers returned error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}
