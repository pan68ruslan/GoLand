package documentstore

import (
	"testing"
)

func newTestCollection() *Collection {
	return &Collection{
		Name:      "users",
		Cfg:       &CollectionConfig{PrimaryKey: "id"},
		Documents: make(map[string]Document),
	}
}

func TestPutAndGet(t *testing.T) {
	c := newTestCollection()
	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "u1"},
			"name": {Type: DocumentFieldTypeString, Value: "Ruslan"},
		},
	}
	if err := c.Put(doc); err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	got, err := c.Get("u1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Fields["name"].Value != "Ruslan" {
		t.Errorf("expected name 'Ruslan', got %v", got.Fields["name"].Value)
	}
	_, err = c.Get("unknown")
	if err == nil {
		t.Errorf("expected error for unknown document, got nil")
	}
}

func TestPutWrongKey(t *testing.T) {
	c := newTestCollection()
	doc := Document{
		Fields: map[string]DocumentField{
			"name": {Type: DocumentFieldTypeString, Value: "Ruslan"},
		},
	}
	if err := c.Put(doc); err == nil {
		t.Errorf("expected error for missing primary key, got nil")
	}
	doc = Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeNumber, Value: 123},
			"name": {Type: DocumentFieldTypeString, Value: "Ruslan"},
		},
	}
	if err := c.Put(doc); err == nil {
		t.Errorf("expected error for wrong key type, got nil")
	}
}

func TestDelete(t *testing.T) {
	c := newTestCollection()
	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "u2"},
			"name": {Type: DocumentFieldTypeString, Value: "Anna"},
		},
	}
	_ = c.Put(doc)
	if e := c.Delete("u2"); e != nil {
		t.Fatalf("Delete failed: %v", e)
	}
	if e := c.Delete("u2"); e == nil {
		t.Errorf("expected error for deleting non-existing document, got nil")
	}
}

func TestList(t *testing.T) {
	c := newTestCollection()
	_, e := c.List()
	if e == nil {
		t.Errorf("expected error for empty collection, got nil")
	}
	_ = c.Put(Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "u3"},
			"name": {Type: DocumentFieldTypeString, Value: "Petro"},
		},
	})
	_ = c.Put(Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "u4"},
			"name": {Type: DocumentFieldTypeString, Value: "Oksana"},
		},
	})
	docs, e := c.List()
	if e != nil {
		t.Fatalf("List failed: %v", e)
	}
	if len(docs) != 2 {
		t.Errorf("expected 2 documents, got %d", len(docs))
	}
}
