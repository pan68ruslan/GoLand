package documentstore

import (
	"encoding/json"
	"fmt"
)

type Collection struct {
	Name      string              `json:"Name"`
	Cfg       *CollectionConfig   `json:"Config"`
	Documents map[string]Document `json:"Documents"`
}

type CollectionConfig struct {
	PrimaryKey string `json:"PrimaryKey"`
}

func (cfg *CollectionConfig) MarshalJSON() ([]byte, error) {
	if cfg == nil {
		return []byte("null"), nil
	}
	out := map[string]interface{}{
		"PrimaryKey": cfg.PrimaryKey,
	}
	return json.Marshal(out)
}

func (cfg *CollectionConfig) UnmarshalJSON(data []byte) error {
	raw := make(map[string]interface{})
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if pk, ok := raw["PrimaryKey"].(string); ok {
		cfg.PrimaryKey = pk
	} else {
		return fmt.Errorf("CollectionConfig: missing or invalid PrimaryKey")
	}
	return nil
}

func (c *Collection) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}
	docs := make(map[string]interface{})
	for key, doc := range c.Documents {
		docs[key] = doc
	}
	out := map[string]interface{}{
		"Name":      c.Name,
		"Config":    c.Cfg,
		"Documents": docs,
	}
	return json.Marshal(out)
}

func (c *Collection) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		StoreLogger.Error("[Collection]UnmarshalJSON() failed")
		return err
	}
	if v, ok := raw["Name"]; ok {
		if err := json.Unmarshal(v, &c.Name); err != nil {
			return err
		}
	}
	if v, ok := raw["Config"]; ok {
		var cfg CollectionConfig
		if err := json.Unmarshal(v, &cfg); err != nil {
			return err
		}
		c.Cfg = &cfg
	}
	if v, ok := raw["Documents"]; ok {
		docs := make(map[string]Document)
		if err := json.Unmarshal(v, &docs); err != nil {
			return err
		}
		c.Documents = docs
	}
	return nil
}

func (c *Collection) Put(doc Document) {
	if c.Cfg == nil || c.Cfg.PrimaryKey == "" {
		StoreLogger.Error("[Collection]CollectionConfig is not configured")
		return
	}
	field, ok := doc.Fields[c.Cfg.PrimaryKey]
	if !ok || field.Type != DocumentFieldTypeString {
		StoreLogger.Error("[Collection]Document is wrong")
		return
	}
	key, ok := field.Value.(string)
	if !ok {
		StoreLogger.Error("[Collection]Key is not a string")
		return
	}
	StoreLogger.Info(fmt.Sprintf("[Collection]The document was added with key '%s'", key))
	c.Documents[key] = doc
}

func (c *Collection) Get(key string) (*Document, bool) {
	doc, ok := c.Documents[key]
	if !ok {
		StoreLogger.Error(fmt.Sprintf("[Collection]The document with key '%s' wasn't found", key))
		return nil, false
	}
	StoreLogger.Info(fmt.Sprintf("[Collection]The document with key '%s' was found", key))
	return &doc, true
}

func (c *Collection) Delete(key string) bool {
	if _, ok := c.Documents[key]; ok {
		StoreLogger.Info(fmt.Sprintf("[Collection]The document with key '%s' was deleted", key))
		delete(c.Documents, key)
		return true
	}
	StoreLogger.Error(fmt.Sprintf("[Collection]The document with key '%s' doesn't exist", key))
	return false
}

func (c *Collection) List() []Document {
	docs := make([]Document, 0, len(c.Documents))
	for _, doc := range c.Documents {
		docs = append(docs, doc)
	}
	if l := len(docs); l < 1 {
		StoreLogger.Error(fmt.Sprintf("[Collection]There are no documents in the collection '%s'", c.Name))
	} else {
		StoreLogger.Info(fmt.Sprintf("[Collection]There are %d documents in the collection '%s'", l, c.Name))
	}
	return docs
}
