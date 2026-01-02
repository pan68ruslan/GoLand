package documentstore

import (
	"encoding/json"
	"errors"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrWrongFieldType  = errors.New("wrong field")
	ErrWrongDocument   = errors.New("wrong Document, type")
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
	Fields map[string]DocumentField `json:"-"`
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
	if d.Fields == nil {
		d.Fields = make(map[string]DocumentField)
	}
	return nil
}

/*func SafetySerialize(v interface{}) ([]byte, error) {
	if doc, ok := v.(Document); ok {
		return json.Marshal(doc)
	}
	return nil, fmt.Errorf("%w: expected Document, got %T", ErrWrongDocument, v)
}*/
