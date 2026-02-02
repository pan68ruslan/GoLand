package documentStore

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType `json:"type"`
	Value interface{}       `json:"value"`
}

type Document struct {
	Fields map[string]DocumentField `json:"fields"`
}

func NewDoc(owner string) *Document {
	return &Document{
		Fields: map[string]DocumentField{
			"id":    DocumentField{Type: DocumentFieldTypeNumber, Value: 0},
			"text":  DocumentField{Type: DocumentFieldTypeString, Value: fmt.Sprintf("Document was created by %s at %s.\n", owner, time.Now().Format("2006-01-02 15:04:05.00"))},
			"owner": DocumentField{Type: DocumentFieldTypeString, Value: owner},
		},
	}
}

func (d *Document) UpdateContent(worker string) error {
	field, ok := d.Fields["text"]
	if !ok {
		slog.Error("field 'text' not found in document")
		return fmt.Errorf("missing 'text' field")
	}
	if field.Type != DocumentFieldTypeString {
		slog.Error("field 'text' has invalid type", "expected", DocumentFieldTypeString, "got", field.Type)
		return fmt.Errorf("'text' field is not a string")
	}
	text, ok := field.Value.(string)
	if !ok {
		slog.Error("field 'text' value is not a string", "value", field.Value)
		return fmt.Errorf("value of 'text' field is not a string")
	}
	updateMsg := fmt.Sprintf("Document was updated by %s at %s.\n", worker, time.Now().Format("2006-01-02 15:04:05.00"))
	text += updateMsg
	field.Value = text
	d.Fields["text"] = field
	slog.Info("document updated successfully", "update", updateMsg)
	return nil
}

// For marshaling
func (d *Document) MarshalJSON() ([]byte, error) {
	out := make(map[string]map[string]interface{})
	for key, field := range d.Fields {
		out[key] = map[string]interface{}{
			"type":  field.Type,
			"value": field.Value,
		}
	}
	return json.Marshal(out)
}

func (d *Document) UnmarshalJSON(data []byte) error {
	fields := make(map[string]DocumentField)
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	d.Fields = fields
	return nil
}
