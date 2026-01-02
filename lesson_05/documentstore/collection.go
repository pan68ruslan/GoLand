package documentstore

import "fmt"

var (
	ErrWrongPrimaryKey     = fmt.Errorf("the config primaryKey is wrong or absent")
	ErrWrongCollectionName = fmt.Errorf("the config collectionName is wrong or empty")
	ErrWrongKeyValue       = fmt.Errorf("the config value is wrong or empty")
	ErrCantFindDocument    = fmt.Errorf("cannot find the document")
	ErrWrongCollection     = fmt.Errorf("the collection is corrupted or empty")
)

type CollectionConfig struct {
	PrimaryKey string `json:"PrimaryKey"`
}

type Collection struct {
	Name      string
	Cfg       *CollectionConfig
	Documents map[string]Document
}

func NewCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if cfg == nil || len(cfg.PrimaryKey) == 0 {
		return nil, ErrWrongPrimaryKey
	}
	if name == "" {
		return nil, ErrWrongCollectionName
	}
	return &Collection{Name: name, Cfg: cfg, Documents: make(map[string]Document)}, nil
}

func (c *Collection) Put(doc Document) error {
	field, ok := doc.Fields[c.Cfg.PrimaryKey]
	if !ok || field.Type != DocumentFieldTypeString {
		return ErrWrongPrimaryKey
	}
	key, ok := field.Value.(string)
	if !ok || len(key) == 0 {
		return ErrWrongKeyValue
	}
	c.Documents[key] = doc
	return nil
}

func (c *Collection) Get(key string) (Document, error) {
	doc, ok := c.Documents[key]
	if ok {
		return doc, nil
	}
	return Document{}, fmt.Errorf("%w, key='%s'", ErrCantFindDocument, key)
}

func (c *Collection) Delete(key string) error {
	if _, ok := c.Documents[key]; ok {
		delete(c.Documents, key)
		return nil
	}
	return fmt.Errorf("%w, key='%s'", ErrCantFindDocument, key)
}

func (c *Collection) List() ([]Document, error) {
	if len(c.Documents) == 0 {
		return nil, ErrWrongCollection
	}
	docs := make([]Document, 0, len(c.Documents))
	for _, doc := range c.Documents {
		docs = append(docs, doc)
	}
	return docs, nil
}
