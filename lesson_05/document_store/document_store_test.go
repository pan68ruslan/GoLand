package Document_Store

import (
	"reflect"
	"testing"
)

func TestMarshalDocument(t *testing.T) {
	doc := Document{
		Fields: map[string]DocumentField{
			"id": {
				Type:  DocumentFieldTypeString,
				Value: "u1",
			},
			"age": {
				Type:  DocumentFieldTypeNumber,
				Value: 30,
			},
			"active": {
				Type:  DocumentFieldTypeBool,
				Value: true,
			},
		},
	}
	result, _ := MarshalDocument(doc)
	expected := map[string]interface{}{
		"id":     "u1",
		"age":    30,
		"active": true,
	}
	if reflect.DeepEqual(result, expected) {
		t.Log("MarshalDocument test Passed")
	} else {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestUnmarshalDocument(t *testing.T) {
	input := map[string]interface{}{
		"id":     "u1",
		"age":    30,
		"active": true,
	}
	doc, err := UnmarshalDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(doc.Fields))
	}
	tests := []struct {
		key      string
		expType  DocumentFieldType
		expValue interface{}
	}{
		{"id", DocumentFieldTypeString, "u1"},
		{"age", DocumentFieldTypeNumber, 30},
		{"active", DocumentFieldTypeBool, true},
	}
	for _, tt := range tests {
		pass := true
		f, ok := doc.Fields[tt.key]
		if !ok {
			pass = false
			t.Fatalf("missing field %s", tt.key)
		}
		if f.Type != tt.expType {
			pass = false
			t.Fatalf("expected type %s, got %s", tt.expType, f.Type)
		}
		if f.Value != tt.expValue {
			pass = false
			t.Fatalf("expected value %v, got %v", tt.expValue, f.Value)
		}
		if pass {
			t.Logf("UnmarshalDocument test Passed")
		}
	}
}

func TestUnmarshalDocumentError(t *testing.T) {
	input := map[string]interface{}{
		"bad": struct{}{}, // unsupported type
	}
	_, err := UnmarshalDocument(input)
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	} else {
		t.Log("TestUnmarshalDocumentError test Passed")
	}
}

func TestCheckFieldType(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected DocumentFieldType
	}{
		{"hello", DocumentFieldTypeString},
		{123, DocumentFieldTypeNumber},
		{3.14, DocumentFieldTypeNumber},
		{true, DocumentFieldTypeBool},
		{[3]int{1, 2, 3}, DocumentFieldTypeArray},
		{map[string]interface{}{"a": 1}, DocumentFieldTypeObject},
	}
	for _, tt := range tests {
		pass := true
		v := reflect.ValueOf(tt.value)
		typ, err := CheckFieldType(v)
		if err != nil {
			pass = false
			t.Fatalf("unexpected error for %v: %v", tt.value, err)
		}
		if typ != tt.expected {
			pass = false
			t.Fatalf("expected %v, got %v for value %v", tt.expected, typ, tt.value)
		}
		if pass {
			t.Logf("UnmarshalDocument test Passed")
		}
	}
}

func TestCheckFieldTypeUnsupported(t *testing.T) {
	type Custom struct{}
	v := reflect.ValueOf(Custom{})
	_, err := CheckFieldType(v)
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	} else {
		t.Log("TestCheckFieldType test Passed")
	}
}
