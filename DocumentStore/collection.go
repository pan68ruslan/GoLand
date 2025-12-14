package DocumentStore

import "fmt"

type Collection struct {
	Name      string
	Cfg       *CollectionConfig
	Documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s *Collection) Put(doc Document) {
	field, ok := doc.Fields[s.Cfg.PrimaryKey]
	if !ok || field.Type != DocumentFieldTypeString {
		fmt.Println("[Collection]PrimaryKey is wrong or absent")
		return
	}
	key, ok := field.Value.(string)
	if !ok {
		fmt.Println("Key is not a string")
		return
	}
	fmt.Printf("[Collection]The document was added with key '%s'\n", key)
	s.Documents[key] = doc
}

func (s *Collection) Get(key string) (*Document, bool) {
	doc, ok := s.Documents[key]
	if !ok {
		fmt.Printf("[Collection]The document with key '%s' wasn't found\n", key)
		return nil, false
	}
	fmt.Printf("[Collection]The document with key '%s' was found\n", key)
	return &doc, true
}

func (s *Collection) Delete(key string) bool {
	if _, ok := s.Documents[key]; ok {
		fmt.Printf("[Collection]The document with key '%s' was deleted\n", key)
		delete(s.Documents, key)
		return true
	}
	fmt.Printf("[Collection]The document with key '%s' doesn't exist\n", key)
	return false
}

func (s *Collection) List() []Document {
	docs := make([]Document, 0)
	for _, doc := range s.Documents {
		docs = append(docs, doc)
	}
	if l := len(docs); l < 1 {
		fmt.Printf("[Collection]There are no documents in the collection '%s'\n", s.Name)
	} else {
		fmt.Printf("[Collection]There are %d documents in the collection '%s'\n", l, s.Name)
		for i, doc := range docs {
			fmt.Printf("[Collection] %d. %s\n", i+1, doc.Fields["title"].Value)
		}
	}
	return docs
}
