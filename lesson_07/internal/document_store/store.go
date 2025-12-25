package document_store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var ErrWrongStore = errors.New("wrong store, type")

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
		return []byte("null"), fmt.Errorf("input collection isn't exists")
	}
	out := make(map[string]json.RawMessage)
	for key, coll := range s.Collections {
		data, err := json.Marshal(coll)
		if err != nil {
			return nil, fmt.Errorf("marshal collection %s: %w", key, err)
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
	if v, ok := raw["Collections"]; ok { // узгоджено з тегом
		colls := make(map[string]json.RawMessage)
		if err := json.Unmarshal(v, &colls); err != nil {
			return err
		}
		if v, ok := raw["Name"]; ok {
			if err := json.Unmarshal(v, &s.Name); err != nil {
				return fmt.Errorf("unmarshal Name: %w", err)
			}
		}
		s.Collections = make(map[string]*Collection)
		for key, collData := range colls {
			var coll Collection
			if err := json.Unmarshal(collData, &coll); err != nil {
				return fmt.Errorf("unmarshal Collections %s: %w", key, err)
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
		fmt.Printf("[Store]The collection '%s' already exists\n", name)
		return false, nil
	}
	col := &Collection{
		Name:      name,
		Cfg:       cfg,
		Documents: make(map[string]Document),
	}
	s.Collections[name] = col
	fmt.Printf("[Store]The collection '%s' was created\n", name)
	return true, col
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	if col, ok := s.Collections[name]; ok {
		fmt.Printf("[Store]The collection '%s' was found\n", name)
		return col, true
	}
	fmt.Printf("[Store]The collection '%s' was not found\n", name)
	return nil, false
}

func (s *Store) DeleteCollection(name string) bool {
	if _, ok := s.Collections[name]; ok {
		fmt.Printf("[Store]The collection '%s' has been deleted\n", name)
		delete(s.Collections, name)
		return true
	}
	fmt.Printf("[Store]The collection '%s' doesn't exist\n", name)
	return false
}

func (s *Store) Dump() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("store is nil")
	}
	return json.MarshalIndent(s, "", "  ")
}

func (s *Store) DumpToFile(filename string) error {
	if s == nil {
		return fmt.Errorf("store is nil")
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal store: %w", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}
	return nil
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	if len(dump) == 0 {
		return nil, fmt.Errorf("empty dump")
	}
	var store Store
	if err := json.Unmarshal(dump, &store); err != nil {
		return nil, fmt.Errorf("cannot unmarshal store: %w", err)
	}
	if store.Collections == nil {
		store.Collections = make(map[string]*Collection)
	}
	return &store, nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}
	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("cannot unmarshal store: %w", err)
	}
	if store.Collections == nil {
		store.Collections = make(map[string]*Collection)
	}
	return &store, nil
}
