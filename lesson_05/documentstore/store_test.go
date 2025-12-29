package documentstore

import (
	"testing"
)

var testStore = "TestDB"
var testCollection = "Users"
var testUser = "User0"
var testId = "u0"

func TestCreateCollection(t *testing.T) {
	store := NewStore(testStore)
	col, err := store.CreateCollection(testCollection, &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if col.Name != testCollection {
		t.Errorf("expected collection name '%s', got %s", testCollection, col.Name)
	}
	_, err = store.CreateCollection(testCollection, &CollectionConfig{PrimaryKey: "id"})
	if err == nil {
		t.Errorf("expected error for duplicate collection, got nil")
	}
	_, err = store.CreateCollection("broken", &CollectionConfig{PrimaryKey: ""})
	if err == nil {
		t.Errorf("expected error for wrong config, got nil")
	}
}

func TestGetCollection(t *testing.T) {
	store := NewStore(testStore)
	_, _ = store.CreateCollection(testCollection, &CollectionConfig{PrimaryKey: "id"})
	col, err := store.GetCollection(testCollection)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if col.Name != testCollection {
		t.Errorf("expected '%s', got %s", testCollection, col.Name)
	}
	_, err = store.GetCollection("unknown")
	if err == nil {
		t.Errorf("expected error for unknown collection, got nil")
	}
}

func TestDeleteCollection(t *testing.T) {
	store := NewStore("TestDB")
	_, _ = store.CreateCollection(testCollection, &CollectionConfig{PrimaryKey: "id"})
	if err := store.DeleteCollection(testCollection); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := store.Collections[testCollection]; ok {
		t.Errorf("expected collection '%s' to be deleted", testCollection)
	}
	if err := store.DeleteCollection(testCollection); err == nil {
		t.Errorf("expected error for deleting non-existing collection, got nil")
	}
}
