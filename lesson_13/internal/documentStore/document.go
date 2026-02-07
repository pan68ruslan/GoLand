package documentStore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
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
	mu     sync.RWMutex
}

func NewDocument(owner string) *Document {
	text := fmt.Sprintf("Document was created by %s at %s.\n", owner, time.Now().Format("2006-01-02 15:04:05.00"))
	return &Document{
		Fields: map[string]DocumentField{
			"id":    DocumentField{Type: DocumentFieldTypeNumber, Value: 0},
			"owner": DocumentField{Type: DocumentFieldTypeString, Value: owner},
			"text":  DocumentField{Type: DocumentFieldTypeString, Value: text},
		},
	}
}

func (d *Document) UpdateContent(owner string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if field, ok := d.Fields["text"]; ok && field.Type == DocumentFieldTypeString {
		if text, ok := field.Value.(string); ok {
			updateMsg := fmt.Sprintf(
				"Document was updated by %s at %s.\n",
				owner,
				time.Now().Format("2006-01-02 15:04:05.00"),
			)
			text += updateMsg
			d.Fields["text"] = DocumentField{
				Type:  DocumentFieldTypeString,
				Value: text,
			}
			slog.Info("document updated successfully",
				"owner", owner,
				"update", updateMsg,
			)
		} else {
			msg := "value of 'text' field is not a string"
			slog.Error(msg, "value", field.Value)
			return errors.New(msg)
		}
	} else {
		msg := "missing 'text' field"
		slog.Error(msg)
		return errors.New(msg)
	}
	return nil
}

// For marshaling
func (d *Document) MarshalJSON() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
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
	d.mu.Lock()
	defer d.mu.Unlock()
	fields := make(map[string]DocumentField)
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	d.Fields = fields
	return nil
}
