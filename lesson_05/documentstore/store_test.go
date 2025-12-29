package documentstore

import (
	"testing"
)

func TestCreateCollection(t *testing.T) {
	store := NewStore("TestDB")
	col, err := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if col.Name != "users" {
		t.Errorf("expected collection name 'users', got %s", col.Name)
	}
	_, err = store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err == nil {
		t.Errorf("expected error for duplicate collection, got nil")
	}
	_, err = store.CreateCollection("broken", &CollectionConfig{PrimaryKey: ""})
	if err == nil {
		t.Errorf("expected error for wrong config, got nil")
	}
}

func TestGetCollection(t *testing.T) {
	store := NewStore("TestDB")
	_, _ = store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	col, err := store.GetCollection("users")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if col.Name != "users" {
		t.Errorf("expected 'users', got %s", col.Name)
	}
	_, err = store.GetCollection("unknown")
	if err == nil {
		t.Errorf("expected error for unknown collection, got nil")
	}
}

func TestDeleteCollection(t *testing.T) {
	store := NewStore("TestDB")
	_, _ = store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err := store.DeleteCollection("users"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := store.Collections["users"]; ok {
		t.Errorf("expected collection 'users' to be deleted")
	}
	if err := store.DeleteCollection("users"); err == nil {
		t.Errorf("expected error for deleting non-existing collection, got nil")
	}
}
