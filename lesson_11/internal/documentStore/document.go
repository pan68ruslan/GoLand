package documentStore

import (
	"encoding/json"
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

func newDoc(name string, field DocumentField) Document {
	return Document{
		Fields: map[string]DocumentField{
			name: {
				Type:  DocumentFieldTypeString,
				Value: field,
			},
		},
	}
}

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
