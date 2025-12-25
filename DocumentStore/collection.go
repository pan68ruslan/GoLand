package DocumentStore

import "fmt"

type Collection struct {
	Name      string
	Cfg       *CollectionConfig
	documents map[string]Document
}

func NewCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if cfg == nil || cfg.PrimaryKey == "" {
		return nil, fmt.Errorf("[Collection]CollectionConfig is not configured")
	}
	return &Collection{
		Name:      name,
		Cfg:       cfg,
		documents: make(map[string]Document),
	}, nil
}

type CollectionConfig struct {
	PrimaryKey string `json:"PrimaryKey"`
}

func (s *Collection) Put(doc Document) {
	if s.Cfg == nil || s.Cfg.PrimaryKey == "" {
		fmt.Println("[Collection]CollectionConfig is not configured")
		return
	}
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
	s.documents[key] = doc
}

func (s *Collection) Get(key string) (*Document, bool) {
	doc, ok := s.documents[key]
	if !ok {
		fmt.Printf("[Collection]The document with key '%s' wasn't found\n", key)
		return nil, false
	}
	fmt.Printf("[Collection]The document with key '%s' was found\n", key)
	return &doc, true
}

func (s *Collection) Delete(key string) bool {
	if _, ok := s.documents[key]; ok {
		fmt.Printf("[Collection]The document with key '%s' was deleted\n", key)
		delete(s.documents, key)
		return true
	}
	fmt.Printf("[Collection]The document with key '%s' doesn't exist\n", key)
	return false
}

func (s *Collection) List() []Document {
	docs := make([]Document, 0)
	for _, doc := range s.documents {
		docs = append(docs, doc)
	}
	if l := len(docs); l < 1 {
		fmt.Printf("[Collection]There are no documents in the collection '%s'\n", s.Name)
	} else {
		fmt.Printf("[Collection]There are %d documents in the collection '%s'\n", l, s.Name)
	}
	return docs
}
