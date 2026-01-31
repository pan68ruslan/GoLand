package documentStore

import (
	"fmt"
	"log/slog"
	"sort"
	"sync"
)

type QueryParams struct {
	Desc     bool
	MinValue *string
	MaxValue *string
}

type CollectionConfig struct {
	PrimaryKey    string   `json:"PrimaryKey"`
	IndexedFields []string `json:"IndexedFields"`
}

type Collection struct {
	Name      string                 `json:"Name"`
	Cfg       *CollectionConfig      `json:"Config"`
	Documents map[string]Document    `json:"Documents"`
	Indexes   map[string]*BinaryTree `json:"Indexes"`
	Logger    *slog.Logger
	mu        sync.RWMutex
}

func NewCollection(name string, cfg *CollectionConfig, logger *slog.Logger) *Collection {
	coll := &Collection{
		Name:      name,
		Cfg:       cfg,
		Documents: make(map[string]Document),
		Indexes:   make(map[string]*BinaryTree),
		Logger:    logger,
	}
	if cfg != nil {
		for _, field := range cfg.IndexedFields {
			coll.Indexes[field] = &BinaryTree{}
		}
	}
	return coll
}

func (c *Collection) PutDocument(doc Document) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Cfg == nil || c.Cfg.PrimaryKey == "" {
		return fmt.Errorf("collection config is not configured")
	}
	field, ok := doc.Fields[c.Cfg.PrimaryKey]
	if !ok || field.Type != DocumentFieldTypeString {
		return fmt.Errorf("document missing primary key field or wrong type")
	}
	key, ok := field.Value.(string)
	if !ok {
		return fmt.Errorf("primary key value is not a string")
	}
	if oldDoc, exists := c.Documents[key]; exists {
		for _, idxField := range c.Cfg.IndexedFields {
			if oldField, ok := oldDoc.Fields[idxField]; ok {
				c.Indexes[idxField].RemoveFromIndex(fmt.Sprintf("%v", oldField.Value), key)
			}
		}
	}
	c.Documents[key] = doc
	for _, idxField := range c.Cfg.IndexedFields {
		if field, ok := doc.Fields[idxField]; ok {
			c.Indexes[idxField].Insert(fmt.Sprintf("%v", field.Value), key)
		}
	}
	return nil
}

func (c *Collection) DeleteDocument(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	doc, exists := c.Documents[key]
	if !exists {
		return false
	}
	for _, idxField := range c.Cfg.IndexedFields {
		if field, ok := doc.Fields[idxField]; ok {
			c.Indexes[idxField].RemoveFromIndex(fmt.Sprintf("%v", field.Value), key)
		}
	}
	delete(c.Documents, key)
	return true
}

func (c *Collection) GetDocument(key string) (*Document, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	doc, ok := c.Documents[key]
	if !ok {
		c.Logger.Info(fmt.Sprintf("[Collection]The document with key '%s' wasn't found", key))
		return nil, false
	}
	c.Logger.Info(fmt.Sprintf("[Collection]The document with key '%s' was found", key))
	return &doc, true
}

func (c *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	idx, exists := c.Indexes[fieldName]
	if !exists {
		return nil, fmt.Errorf("index does not exist for field %s", fieldName)
	}
	if params.MinValue == nil || params.MaxValue == nil || *params.MinValue > *params.MaxValue {
		return nil, fmt.Errorf("parameters do not exist or wrong")
	}
	ids := idx.RangeSearch(params.MinValue, params.MaxValue)
	if params.Desc {
		sort.Slice(ids, func(i, j int) bool { return ids[i] > ids[j] })
	} else {
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	}
	var result []Document
	for _, id := range ids {
		if doc, ok := c.Documents[id]; ok {
			result = append(result, doc)
		}
	}
	c.Logger.Info("Query executed",
		slog.String("collection", c.Name),
		slog.String("field", fieldName),
		slog.Int("found", len(result)))
	return result, nil
}
