package documentstore

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrMarshal   = errors.New("MarshalDocument: expected struct")
	ErrUnmarshal = errors.New("UnmarshalDocument: input doc must be pointer to struct")
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

func MarshalStructureToDocument(v interface{}) (Document, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	if k := val.Kind(); k != reflect.Struct {
		return Document{}, fmt.Errorf("%w, got %s", ErrMarshal, k)
	}
	fields := make(map[string]DocumentField)
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		key := fieldType.Tag.Get("json")
		if key == "" {
			key = fieldType.Name
		}
		var docType DocumentFieldType
		switch fieldVal.Kind() {
		case reflect.String:
			docType = DocumentFieldTypeString
		case reflect.Int, reflect.Int64, reflect.Float64:
			docType = DocumentFieldTypeNumber
		case reflect.Bool:
			docType = DocumentFieldTypeBool
		case reflect.Slice, reflect.Array:
			docType = DocumentFieldTypeArray
		case reflect.Struct, reflect.Map:
			docType = DocumentFieldTypeObject
		default:
			docType = DocumentFieldTypeString
		}
		fields[key] = DocumentField{
			Type:  docType,
			Value: fieldVal.Interface(),
		}
	}
	return Document{Fields: fields}, nil
}

func UnmarshalDocumentToStructure(doc *Document, out interface{}) error {
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return ErrUnmarshal
	}
	structVal := val.Elem()
	structType := structVal.Type()
	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)
		key := fieldType.Tag.Get("json")
		if key == "" {
			key = fieldType.Name
		} else {
			if idx := strings.Index(key, ","); idx != -1 {
				key = key[:idx]
			}
		}
		docField, ok := doc.Fields[key]
		if !ok {
			continue
		}
		if field.CanSet() {
			fieldVal := reflect.ValueOf(docField.Value)
			if fieldVal.Type().AssignableTo(field.Type()) {
				field.Set(fieldVal)
			} else if fieldVal.Type().ConvertibleTo(field.Type()) {
				field.Set(fieldVal.Convert(field.Type()))
			}
		}
	}
	return nil
}
