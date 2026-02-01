package documentStore

import (
	"fmt"
	"log/slog"
	"math"
	"sort"
	"sync"
)

type QueryParams struct {
	Desc     bool
	MinValue int
	MaxValue int
}

func NewQueryParams() *QueryParams {
	return &QueryParams{
		Desc:     false,
		MinValue: -1,
		MaxValue: math.MaxInt32,
	}
}

type CollectionConfig struct {
	PrimaryKey    string   `json:"PrimaryKey"`
	IndexedFields []string `json:"Indexes"`
}

func NewConfig() *CollectionConfig {
	key := "id"
	return &CollectionConfig{
		PrimaryKey:    key,
		IndexedFields: []string{key}, //{"id", "name"},
	}
}

type Collection struct {
	Name      string                 `json:"Name"`
	Cfg       *CollectionConfig      `json:"Config"`
	Documents map[int]Document       `json:"Documents"`
	Indexes   map[string]*BinaryTree `json:"Indexes"`
	Logger    *slog.Logger
	mu        sync.RWMutex
}

func NewCollection(name string, logger *slog.Logger) Collection {
	cfg := NewConfig()
	coll := Collection{
		Name:      name,
		Cfg:       cfg,
		Documents: make(map[int]Document),
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
	/*if c.Cfg == nil || c.Cfg.PrimaryKey == "" {
		return fmt.Errorf("collection config is not configured")
	}*/
	field, ok := doc.Fields[c.Cfg.PrimaryKey]
	if !ok || field.Type != DocumentFieldTypeNumber {
		return fmt.Errorf("document missing primary key field or wrong type")
	}
	key, ok := field.Value.(int)
	if !ok {
		return fmt.Errorf("primary key value is not a number")
	}
	if oldDoc, exists := c.Documents[key]; exists {
		for _, idxField := range c.Cfg.IndexedFields {
			if oldField, ok := oldDoc.Fields[idxField]; ok {
				c.Indexes[idxField].RemoveFromIndex(oldField.Value.(int), key)
			}
		}
	}
	c.Documents[key] = doc
	for _, idxField := range c.Cfg.IndexedFields {
		if field, ok := doc.Fields[idxField]; ok {
			c.Indexes[idxField].Insert(field.Value.(int), key)
		}
	}
	return nil
}

func (c *Collection) DeleteDocument(key int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	doc, exists := c.Documents[key]
	if !exists {
		return false
	}
	for _, idxField := range c.Cfg.IndexedFields {
		if field, ok := doc.Fields[idxField]; ok {
			c.Indexes[idxField].RemoveFromIndex(field.Value.(int), key)
		}
	}
	delete(c.Documents, key)
	return true
}

func (c *Collection) GetDocument(key int) (*Document, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	doc, ok := c.Documents[key]
	if !ok {
		c.Logger.Info(fmt.Sprintf("[Collection]The document with key '%d' wasn't found", key))
		return nil, false
	}
	c.Logger.Info(fmt.Sprintf("[Collection]The document with key '%d' was found", key))
	return &doc, true
}

func (c *Collection) Query( /*fieldName string, */ params QueryParams) ([]Document, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	idx, exists := c.Indexes[c.Cfg.PrimaryKey] //c.Indexes[fieldName]
	if !exists {
		return nil, fmt.Errorf("index does not exist for field %s", c.Cfg.PrimaryKey)
	}
	if params.MinValue > params.MaxValue {
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
		slog.String("field", c.Cfg.PrimaryKey),
		slog.Int("found", len(result)))
	return result, nil
}

func (c *Collection) MaxId() int {
	max := 0
	for _, doc := range c.Documents {
		if m := doc.Fields["id"].Value.(int); m > max {
			max = m
		}
	}
	return max
}
