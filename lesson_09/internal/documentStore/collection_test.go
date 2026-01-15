package documentStore

import (
	"encoding/json"
	"testing"
)

func newDoc(key string, title string) Document {
	return Document{
		Fields: map[string]DocumentField{
			"key":   {Type: DocumentFieldTypeString, Value: key},
			"title": {Type: DocumentFieldTypeString, Value: title},
		},
	}
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
	c := &Collection{
		Name:      "TestCollection",
		Cfg:       &CollectionConfig{PrimaryKey: "key"},
		Documents: map[string]Document{"doc1": newDoc("doc1", "first")},
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

func TestCollection_PutGetDeleteList(t *testing.T) {
	c := &Collection{
		Name:      "Docs",
		Cfg:       &CollectionConfig{PrimaryKey: "key"},
		Documents: make(map[string]Document),
	}
	doc := newDoc("doc1", "first")
	c.Put(doc)
	if len(c.Documents) != 1 {
		t.Errorf("expected 1 document, got %d", len(c.Documents))
	}
	got, ok := c.Get("doc1")
	if !ok {
		t.Errorf("expected to find doc1")
	}
	if got != nil && got.Fields["title"].Value != "first" {
		t.Errorf("expected title 'first', got %v", got.Fields["title"].Value)
	}
	if !c.Delete("doc1") {
		t.Errorf("expected delete success")
	}
	if len(c.Documents) != 0 {
		t.Errorf("expected 0 documents after delete, got %d", len(c.Documents))
	}
	c.Put(newDoc("doc2", "second"))
	list := c.List()
	if len(list) != 1 {
		t.Errorf("expected 1 document in list, got %d", len(list))
	}
}
