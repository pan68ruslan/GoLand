package documentStore

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
)

func newTestCollection() *Collection {
	cfg := &CollectionConfig{
		PrimaryKey:    "id",
		IndexedFields: []string{"id", "name"},
	}
	return NewCollection("TestCollection", cfg, slog.Default())
}

func TestCollection_AddGetDelete(t *testing.T) {
	coll := newTestCollection()
	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "u1"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}
	if err := coll.AddDocument(doc); err != nil {
		t.Fatalf("AddDocument failed: %v", err)
	}
	got, ok := coll.GetDocument("u1")
	if !ok {
		t.Fatalf("Get failed: document not found")
	}
	if got.Fields["name"].Value != "Alice" {
		t.Errorf("expected name=Alice, got %v", got.Fields["name"].Value)
	}
	if !coll.DeleteDocument("u1") {
		t.Fatalf("DeleteDocument failed: document not deleted")
	}
	if _, ok := coll.GetDocument("u1"); ok {
		t.Errorf("expected document to be deleted")
	}
}

func TestListDocuments(t *testing.T) {
	coll := newTestCollection()
	for i := 0; i < 3; i++ {
		doc := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: fmt.Sprintf("u%d", i)},
				"name": {Type: DocumentFieldTypeString, Value: fmt.Sprintf("User%d", i)},
			},
		}
		_ = coll.AddDocument(doc)
	}
	docs := coll.ListDocuments()
	if len(docs) != 3 {
		t.Errorf("expected 3 documents, got %d", len(docs))
	}
}

func TestCollection_Query(t *testing.T) {
	coll := newTestCollection()
	docs := []Document{
		{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "u1"}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}},
		{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "u2"}, "name": {Type: DocumentFieldTypeString, Value: "Bob"}}},
		{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "u3"}, "name": {Type: DocumentFieldTypeString, Value: "Charlie"}}},
	}
	for _, d := range docs {
		if err := coll.AddDocument(d); err != nil {
			t.Fatalf("AddDocument failed: %v", err)
		}
	}
	minParam := "A"
	maxParam := "C"
	params := QueryParams{MinValue: &minParam, MaxValue: &maxParam, Desc: false}
	result, err := coll.Query("name", params)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
	names := []string{result[0].Fields["name"].Value.(string), result[1].Fields["name"].Value.(string)}
	if !(contains(names, "Alice") && contains(names, "Bob")) {
		t.Errorf("expected Alice and Bob in results, got %v", names)
	}
}

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func TestCollectionConfig_MarshalUnmarshal(t *testing.T) {
	cfg := &CollectionConfig{PrimaryKey: "key"}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var cfg2 CollectionConfig
	if err := json.Unmarshal(data, &cfg2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if cfg2.PrimaryKey != "key" {
		t.Errorf("expected PrimaryKey 'key', got %s", cfg2.PrimaryKey)
	}
}

func TestCollection_MarshalUnmarshal(t *testing.T) {
	c := newTestCollection()
	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "id"},
			"name": {Type: DocumentFieldTypeString, Value: "name"},
		},
	}
	_ = c.AddDocument(doc)
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var c2 Collection
	if err := json.Unmarshal(data, &c2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if c2.Name != "TestCollection" {
		t.Errorf("expected Name 'TestCollection', got %s", c2.Name)
	}
	if len(c2.Documents) != 1 {
		t.Errorf("expected 1 document, got %d", len(c2.Documents))
	}
}
