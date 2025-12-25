package users

import (
	"testing"
)

func newTestService() *Service {
	coll := &Collection{
		Title: "Test",
		Users: make(map[string]*User),
	}
	return &Service{coll: coll}
}

func TestCreateUser(t *testing.T) {
	svc := newTestService()
	u, err := svc.CreateUser("u1", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	switch {
	case u == nil || u.ID != "u1" || u.Name != "Alice":
		t.Fatalf("unexpected user data: %+v", u)
	}
	_, err = svc.CreateUser("u1", "Alice")
	if err == nil {
		t.Fatalf("expected error for duplicate user, got nil")
	}
}

func TestGetUser(t *testing.T) {
	svc := newTestService()
	svc.CreateUser("u1", "Alice")
	u, err := svc.GetUser("u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil || u.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", u.Name)
	}
	_, err = svc.GetUser("u2")
	if err == nil {
		t.Fatalf("expected error for missing user, got nil")
	}
}

func TestListUsers(t *testing.T) {
	svc := newTestService()
	svc.CreateUser("u1", "Alice")
	svc.CreateUser("u2", "Bob")
	list, err := svc.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 users, got %d", len(list))
	}
}

func TestDeleteUser(t *testing.T) {
	svc := newTestService()
	svc.CreateUser("u1", "Alice")
	u, err := svc.GetUser("u1")
	if u == nil || err != nil {
		t.Fatalf("unexpected error for deleted user u1, %v", err)
	}
	err = svc.DeleteUser("u2")
	if err == nil {
		t.Fatalf("expected error for existing user u2")
	}
	err = svc.DeleteUser("u1")
	if err != nil {
		t.Fatalf("unexpected error for missing user u1, %v", err)
	}
}
