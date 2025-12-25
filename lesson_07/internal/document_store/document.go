package document_store

import (
	"encoding/json"
	"errors"
	"fmt"
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

/*func (d Document) UnmarshalJSON(data []byte) error {
	raw := make(map[string]map[string]interface{})
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	fields := make(map[string]DocumentField)
	for key, val := range raw {
		t, ok := val["type"].(string)
		if !ok {
			return fmt.Errorf("%w: missing type for field %s", ErrWrongFieldType, key)
		}
		dfType := DocumentFieldType(t)
		value := val["value"]
		switch dfType {
		case DocumentFieldTypeString:
			if _, ok := value.(string); !ok {
				return fmt.Errorf("%w: expected string for %s", ErrWrongFieldType, key)
			}
		case DocumentFieldTypeNumber:
			if _, ok := value.(float64); !ok {
				return fmt.Errorf("%w: expected number for %s", ErrWrongFieldType, key)
			}
		case DocumentFieldTypeBool:
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("%w: expected bool for %s", ErrWrongFieldType, key)
			}
		case DocumentFieldTypeArray:
			if _, ok := value.([]interface{}); !ok {
				return fmt.Errorf("%w: expected array for %s", ErrWrongFieldType, key)
			}
		case DocumentFieldTypeObject:
			if _, ok := value.(map[string]interface{}); !ok {
				return fmt.Errorf("%w: expected object for %s", ErrWrongFieldType, key)
			}
		default:
			return fmt.Errorf("%w: %s", ErrUnsupportedType, dfType)
		}
		fields[key] = DocumentField{
			Type:  dfType,
			Value: value,
		}
	}
	d.Fields = fields
	return nil
}*/

func SafetySerialize(v interface{}) ([]byte, error) {
	if doc, ok := v.(Document); ok {
		return json.Marshal(doc)
	}
	return nil, fmt.Errorf("%w: expected Document, got %T", ErrWrongDocument, v)
}
