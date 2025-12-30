package documentstore

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

var StoreLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

type Store struct {
	Name        string
	Collections map[string]*Collection `json:"Collections"`
}

func NewStore(name string) *Store {
	var nm string
	if name != "" {
		nm = name
	} else {
		nm = "NewStore"
	}
	return &Store{
		Name:        nm,
		Collections: make(map[string]*Collection),
	}
}

func (s *Store) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), fmt.Errorf("[Store]Input collection isn't exists")
	}
	out := make(map[string]json.RawMessage)
	for key, coll := range s.Collections {
		data, err := json.Marshal(coll)
		if err != nil {
			return nil, fmt.Errorf("[Store]Marshal collection %s: %w", key, err)
		}
		out[key] = data
	}
	return json.Marshal(map[string]interface{}{
		"Name":        s.Name,
		"Collections": out,
	})
}

func (s *Store) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["Collections"]; ok {
		colls := make(map[string]json.RawMessage)
		if err := json.Unmarshal(v, &colls); err != nil {
			return err
		}
		if v, ok := raw["Name"]; ok {
			if err := json.Unmarshal(v, &s.Name); err != nil {
				return fmt.Errorf("[Store]Unmarshal Name: %w", err)
			}
		}
		s.Collections = make(map[string]*Collection)
		for key, collData := range colls {
			var coll Collection
			if err := json.Unmarshal(collData, &coll); err != nil {
				return fmt.Errorf("[Store]Unmarshal Collections %s: %w", key, err)
			}
			s.Collections[key] = &coll
		}
	} else {
		s.Collections = make(map[string]*Collection)
	}
	return nil
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	if _, ok := s.Collections[name]; ok {
		StoreLogger.Error(fmt.Sprintf("[Store]The collection '%s' already exists", name))
		return false, nil
	}
	col := &Collection{
		Name:      name,
		Cfg:       cfg,
		Documents: make(map[string]Document),
	}
	s.Collections[name] = col
	StoreLogger.Info(fmt.Sprintf("[Store]The collection '%s' was created", name))
	return true, col
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	if col, ok := s.Collections[name]; ok {
		StoreLogger.Info(fmt.Sprintf("[Store]The collection '%s' was found", name))
		return col, true
	}
	StoreLogger.Error(fmt.Sprintf("[Store]The collection '%s' was not found", name))
	return nil, false
}

func (s *Store) DeleteCollection(name string) bool {
	if _, ok := s.Collections[name]; ok {
		StoreLogger.Info(fmt.Sprintf("[Store]The collection '%s' has been deleted", name))
		delete(s.Collections, name)
		return true
	}
	StoreLogger.Error(fmt.Sprintf("[Store]The collection '%s' doesn't exist", name))
	return false
}

func (s *Store) Dump() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("[Store]The store is nil")
	}
	return json.MarshalIndent(s, "", "  ")
}

func (s *Store) DumpToFile(filename string) error {
	if s == nil {
		return fmt.Errorf("[Store]The store is nil")
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("[Store] Cannot marshal store: %w", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("[Store] Cannot write file: %w", err)
	}
	return nil
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	if len(dump) == 0 {
		return nil, fmt.Errorf("[Store]The dump is empty")
	}
	var store Store
	if err := json.Unmarshal(dump, &store); err != nil {
		return nil, fmt.Errorf("[Store]Cannot unmarshal store: %w", err)
	}
	if store.Collections == nil {
		store.Collections = make(map[string]*Collection)
	}
	return &store, nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("[Store]Cannot read file: %w", err)
	}
	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("[Store]Cannot unmarshal store: %w", err)
	}
	if store.Collections == nil {
		store.Collections = make(map[string]*Collection)
	}
	return &store, nil
}
