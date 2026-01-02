package documentstore

import (
	"testing"
)

type TestStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestMarshalStructureToDocument(t *testing.T) {
	original := TestStruct{ID: testId, Name: testUser}

	doc, err := MarshalStructureToDocument(original)
	if err != nil {
		t.Fatalf("MarshalStructureToDocument failed: %v", err)
	}
	if doc.Fields["id"].Value != testId {
		t.Errorf("expected id='%s', got %v", testId, doc.Fields["id"].Value)
	}
	if doc.Fields["name"].Value != testUser {
		t.Errorf("expected name=%s, got %v", testUser, doc.Fields["name"].Value)
	}
}

func TestMarshalStructureToDocumentError(t *testing.T) {
	_, err := MarshalStructureToDocument(123)
	if err == nil {
		t.Errorf("expected error for non-struct input, got nil")
	}
}

func TestUnmarshalDocumentToStructure(t *testing.T) {
	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: testId},
			"name": {Type: DocumentFieldTypeString, Value: testUser},
		},
	}
	var result TestStruct
	if err := UnmarshalDocumentToStructure(&doc, &result); err != nil {
		t.Fatalf("UnmarshalDocumentToStructure failed: %v", err)
	}
	if result.ID != testId || result.Name != testUser {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestUnmarshalDocumentToStructureError(t *testing.T) {
	doc := Document{Fields: map[string]DocumentField{}}
	var notStruct int
	err := UnmarshalDocumentToStructure(&doc, &notStruct)
	if err == nil {
		t.Errorf("expected error for non-struct out, got nil")
	}
}

func TestUnmarshalDocumentMissingField(t *testing.T) {
	doc := Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: testId},
		},
	}
	var result TestStruct
	if err := UnmarshalDocumentToStructure(&doc, &result); err != nil {
		t.Fatalf("UnmarshalDocumentToStructure failed: %v", err)
	}
	if result.Name != "" {
		t.Errorf("expected Name='', got '%s'", result.Name)
	}
}
