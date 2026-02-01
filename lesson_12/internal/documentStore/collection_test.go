package documentStore

import (
	"encoding/json"
	"log/slog"
	"testing"
)

func TestCollection_AddGetDelete(t *testing.T) {
	coll := NewCollection("users", slog.Default())
	id := 1
	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeNumber, Value: id},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}
	if err := coll.PutDocument(doc); err != nil {
		t.Fatalf("AddDocument failed: %v", err)
	}

	got, ok := coll.GetDocument(id)
	if !ok {
		t.Fatalf("Get failed: document not found")
	}
	if got.Fields["name"].Value != "Alice" {
		t.Errorf("expected name=Alice, got %v", got.Fields["name"].Value)
	}

	if !coll.DeleteDocument(id) {
		t.Fatalf("DeleteDocument failed: document not deleted")
	}
	if _, ok := coll.GetDocument(id); ok {
		t.Errorf("expected document to be deleted")
	}
}

func TestCollection_Query(t *testing.T) {
	coll := NewCollection("users", slog.Default())
	id1 := 1
	id2 := 2
	id3 := 3
	docs := []Document{
		{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeNumber, Value: id1}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}},
		{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeNumber, Value: id2}, "name": {Type: DocumentFieldTypeString, Value: "Bob"}}},
		{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeNumber, Value: id3}, "name": {Type: DocumentFieldTypeString, Value: "Charlie"}}},
	}
	for _, d := range docs {
		if err := coll.PutDocument(d); err != nil {
			t.Fatalf("AddDocument failed: %v", err)
		}
	}
	minParam := id1
	maxParam := id3
	params := QueryParams{MinValue: minParam, MaxValue: maxParam, Desc: false}
	result, err := coll.Query(params)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(result) != 3 {
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
	cfg := NewConfig()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var cfg2 CollectionConfig
	if err := json.Unmarshal(data, &cfg2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if cfg2.PrimaryKey != "id" {
		t.Errorf("expected PrimaryKey 'key', got %s", cfg2.PrimaryKey)
	}
}

func TestCollection_MarshalUnmarshal(t *testing.T) {
	id := 1
	c := &Collection{
		Name:      "TestCollection",
		Cfg:       &CollectionConfig{PrimaryKey: "id"},
		Documents: map[int]Document{id: NewDoc("user")},
	}
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
