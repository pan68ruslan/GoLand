package DocStore

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

var documents = map[string]*Document{}

func Put(doc *Document) {
	// 1. Перевірити, що документ має поле "key" типу string
	field, ok := doc.Fields["key"]
	if !ok || field.Type != DocumentFieldTypeString {
		return // якщо немає ключа або він не string — нічого не робимо
	}
	// 2. Отримати значення ключа
	key, ok := field.Value.(string)
	if !ok {
		return
	}
	// 3. Додати документ у мапу
	documents[key] = doc
}

func Get(key string) (*Document, bool) {
	doc, ok := documents[key]
	if !ok {
		return nil, false
	}
	return doc, true
}

func Delete(key string) bool {
	if _, ok := documents[key]; ok {
		delete(documents, key)
		return true
	}
	return false
}

func List() []*Document {
	result := make([]*Document, 0, len(documents))
	for _, doc := range documents {
		result = append(result, doc)
	}
	return result
}
