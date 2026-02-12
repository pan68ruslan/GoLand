package documentStore

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestUpdateContentSuccess(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"text": {Type: DocumentFieldTypeString, Value: "Initial content. "},
		},
	}
	err := doc.UpdateContent("Ruslan")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	updated := doc.Fields["text"].Value.(string)
	if !strings.Contains(updated, "Document was updated by Ruslan") {
		t.Errorf("expected update message, got %s", updated)
	}
}

func TestUpdateContentMissingTextField(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"title": {Type: DocumentFieldTypeString, Value: "No text field"},
		},
	}
	err := doc.UpdateContent("Ruslan")
	if err == nil {
		t.Fatalf("expected error for missing text field")
	}
	if !strings.Contains(err.Error(), "missing 'text' field") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUpdateContentWrongType(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"text": {Type: DocumentFieldTypeString, Value: 123}, // not a string
		},
	}
	err := doc.UpdateContent("Ruslan")
	if err == nil {
		t.Fatalf("expected error for wrong type")
	}
	if !errors.Is(err, errors.New("value of 'text' field is not a string")) {
		if !strings.Contains(err.Error(), "value of 'text' field is not a string") {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestUpdateContentTimestamp(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"text": {Type: DocumentFieldTypeString, Value: "Initial"},
		},
	}
	err := doc.UpdateContent("Ruslan")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	updated := doc.Fields["text"].Value.(string)
	now := time.Now().Format("2006-01-02")
	if !strings.Contains(updated, now) {
		t.Errorf("expected timestamp in update message, got %s", updated)
	}
}

func TestDocument_MarshalUnmarshal_StringField(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"title": {Type: DocumentFieldTypeString, Value: "MyTitle"},
			"text":  {Type: DocumentFieldTypeString, Value: "MyText"},
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
	field, ok := doc2.Fields["title"]
	if !ok {
		t.Errorf("expected field 'title' present")
	}
	if field.Type != DocumentFieldTypeString {
		t.Errorf("expected type 'string', got %s", field.Type)
	}
	if field.Value != "MyTitle" {
		t.Errorf("expected value 'MyTitle', got %v", field.Value)
	}
}

func TestDocument_MarshalUnmarshal_NumberAndBool(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"pages":      {Type: DocumentFieldTypeNumber, Value: 42},
			"isApproved": {Type: DocumentFieldTypeBool, Value: true},
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
	if v := doc2.Fields["pages"].Value; v != float64(42) {
		t.Errorf("expected 42, got %v", v)
	}
	if v := doc2.Fields["isApproved"].Value; v != true {
		t.Errorf("expected true, got %v", v)
	}
}

func TestDocument_UnmarshalEmpty(t *testing.T) {
	var doc Document
	err := json.Unmarshal([]byte(`{}`), &doc)
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
	if string(data) != "{}" {
		t.Errorf("expected '{}', got %s", string(data))
	}
}
