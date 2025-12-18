package Document_Store

import (
	"errors"
	"fmt"
	"reflect"
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
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

func MarshalDocument(inputDoc interface{}) (map[string]interface{}, error) {
	t := reflect.TypeOf(inputDoc)
	if t != reflect.TypeOf(Document{}) {
		return nil, fmt.Errorf("MarshalDocument with %w, type=%v", ErrWrongDocument, t)
	}
	result := make(map[string]interface{})
	v := reflect.ValueOf(inputDoc)
	fields := v.FieldByName("Fields")
	for _, key := range fields.MapKeys() {
		df := fields.MapIndex(key).Interface().(DocumentField)
		result[key.String()] = df.Value
	}
	return result, nil
}

func CheckFieldType(v reflect.Value) (DocumentFieldType, error) {
	var k = v.Kind()
	switch k {
	case reflect.String:
		return DocumentFieldTypeString, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return DocumentFieldTypeNumber, nil
	case reflect.Bool:
		return DocumentFieldTypeBool, nil
	case reflect.Array:
		return DocumentFieldTypeArray, nil
	case reflect.Map:
		return DocumentFieldTypeObject, nil
	default:
		return "", fmt.Errorf("CheckFieldType with %w: type=%v", ErrUnsupportedType, k)
	}
}

func UnmarshalDocument(inputMap map[string]interface{}) (Document, error) {
	outputDoc := Document{
		Fields: make(map[string]DocumentField),
	}
	for key, value := range inputMap {
		v := reflect.ValueOf(value)
		fieldType, err := CheckFieldType(v)
		if err != nil {
			return Document{}, fmt.Errorf("UnmarshalDocument with %w: key=%s, value=%v, err=%v", ErrWrongFieldType, key, value, err)
		}
		outputDoc.Fields[key] = DocumentField{
			Type:  fieldType,
			Value: value,
		}
	}
	return outputDoc, nil
}
