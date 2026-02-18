package documentStore

import (
	"encoding/json"
	"testing"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument("User")
	if doc == nil {
		t.Fatal("NewDocument повернув nil")
	}
	if doc.Fields["owner"].Value != "User" {
		t.Errorf("очікував owner=User, отримав %v", doc.Fields["owner"].Value)
	}
	if _, ok := doc.Fields["text"]; !ok {
		t.Error("очікував поле 'text'")
	}
}

func TestMarshalJSON(t *testing.T) {
	doc := NewDocument("User")
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("MarshalJSON fails: %v", err)
	}
	if len(data) == 0 {
		t.Error("JSON is expected")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	doc := NewDocument("User")
	data, _ := json.Marshal(doc)
	var newDoc Document
	err := json.Unmarshal(data, &newDoc)
	if err != nil {
		t.Fatalf("UnmarshalJSON повернув помилку: %v", err)
	}
	if newDoc.Fields["owner"].Value != "User" {
		t.Errorf("expected owner=User, got %v", newDoc.Fields["owner"].Value)
	}
}

func TestDocument_MarshalUnmarshal_NumberAndString(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"id":    DocumentField{Type: DocumentFieldTypeNumber, Value: 0},
			"owner": DocumentField{Type: DocumentFieldTypeString, Value: "owner"},
			"text":  DocumentField{Type: DocumentFieldTypeString, Value: "text"},
		},
	}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var doc2 Document
	if err := json.Unmarshal(data, &doc2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if v, ok := doc2.Fields["id"].Value.(float64); !ok || v != 0 {
		t.Errorf("expected id=0, got %v", doc2.Fields["id"].Value)
	}
	if v, ok := doc2.Fields["owner"].Value.(string); !ok || v != "owner" {
		t.Errorf("expected text='owner', got %v", doc2.Fields["owner"].Value)
	}
	if v, ok := doc2.Fields["text"].Value.(string); !ok || v != "text" {
		t.Errorf("expected text='text', got %v", doc2.Fields["text"].Value)
	}
}

func TestDocument_UnmarshalEmpty(t *testing.T) {
	var doc Document
	err := json.Unmarshal([]byte(`{"fields":{}}`), &doc)
	if err != nil {
		t.Fatalf("unmarshal empty failed: %v", err)
	}
	if doc.Fields == nil {
		t.Errorf("expected Fields initialized, got nil")
	}
	if len(doc.Fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(doc.Fields))
	}
}

func TestDocument_MarshalEmpty(t *testing.T) {
	doc := &Document{Fields: make(map[string]DocumentField)}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal empty failed: %v", err)
	}
	if string(data) != "{\"fields\":{}}" {
		t.Errorf("expected '{}', got %s", string(data))
	}
}
