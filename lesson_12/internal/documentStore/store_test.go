package documentStore

import (
	"log/slog"
	"os"
	"testing"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func newCfg() *CollectionConfig {
	return &CollectionConfig{PrimaryKey: "key"}
}

func TestNewStore(t *testing.T) {
	s := NewStore("MyStore", logger)
	if s.Name != "MyStore" {
		t.Errorf("expected Name 'MyStore', got %s", s.Name)
	}
	if s.Collections == nil {
		t.Errorf("expected Collections map initialized")
	}
}

func TestCreateGetDeleteCollection(t *testing.T) {
	s := NewStore("TestStore", logger)
	ok, coll := s.CreateCollection("Coll1", logger)
	if !ok || coll == nil {
		t.Fatalf("expected collection created")
	}
	if _, exists := s.Collections["Coll1"]; !exists {
		t.Errorf("collection not stored in map")
	}
	got, ok := s.GetCollection("Coll1")
	if !ok || got == nil {
		t.Errorf("expected to get existing collection")
	}
	if !s.DeleteCollection("Coll1") {
		t.Errorf("expected delete success")
	}
	if _, ok := s.Collections["Coll1"]; ok {
		t.Errorf("collection should be deleted")
	}
	if s.DeleteCollection("CollX") {
		t.Errorf("expected delete non-existent to return false")
	}
}

func TestMarshalUnmarshalStore(t *testing.T) {
	s := NewStore("MarshalStore", logger)
	s.CreateCollection("Coll1", logger)
	data, err := s.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var s2 Store
	if err := s2.UnmarshalJSON(data); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if s2.Name != "MarshalStore" {
		t.Errorf("expected Name 'MarshalStore', got %s", s2.Name)
	}
	if len(s2.Collections) != 1 {
		t.Errorf("expected 1 collection, got %d", len(s2.Collections))
	}
}

func TestDumpAndNewStoreFromDump(t *testing.T) {
	s := NewStore("DumpStore", logger)
	s.CreateCollection("Coll1", logger)
	data, err := s.Dump()
	if err != nil {
		t.Fatalf("dump failed: %v", err)
	}
	s2, err := NewStoreFromDump(data)
	if err != nil {
		t.Fatalf("NewStoreFromDump failed: %v", err)
	}
	if s2.Name != s.Name+"Restored" {
		t.Errorf("expected Name 'DumpStore', got %s", s2.Name)
	}
	if len(s2.Collections) != 1 {
		t.Errorf("expected 1 collection, got %d", len(s2.Collections))
	}
}

func TestDumpToFileAndNewStoreFromFile(t *testing.T) {
	s := NewStore("FileStore", logger)
	s.CreateCollection("Coll1", logger)
	tmpfile := "test_store.json"
	defer os.Remove(tmpfile)
	if err := s.DumpToFile(tmpfile); err != nil {
		t.Fatalf("DumpToFile failed: %v", err)
	}
	s2, err := NewStoreFromFile(tmpfile)
	if err != nil {
		t.Fatalf("NewStoreFromFile failed: %v", err)
	}
	if s2.Name != s.Name+"Restored" {
		t.Errorf("expected Name 'FileStore', got %s", s2.Name)
	}
	if len(s2.Collections) != 1 {
		t.Errorf("expected 1 collection, got %d", len(s2.Collections))
	}
}
