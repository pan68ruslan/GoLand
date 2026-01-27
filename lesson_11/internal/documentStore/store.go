package documentStore

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type Store struct {
	Name        string
	Collections map[string]*Collection `json:"Collections"`
	Logger      *slog.Logger
}

func NewStore(name string, logger *slog.Logger) *Store {
	var nm string
	if name != "" {
		nm = name
	} else {
		nm = "NewStore"
	}
	return &Store{
		Name:        nm,
		Collections: make(map[string]*Collection),
		Logger:      logger,
	}
}

func (s *Store) AddCollection(name string, cfg *CollectionConfig) {
	s.Collections[name] = NewCollection(name, cfg, s.Logger)
	s.Logger.Info("Added collection",
		slog.String("store", s.Name),
		slog.String("collection", name))
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig, logger *slog.Logger) (bool, *Collection) {
	if _, ok := s.Collections[name]; ok {
		fmt.Printf("[Store]The collection '%s' already exists\n", name)
		return false, nil
	}
	col := &Collection{
		Name:      name,
		Cfg:       cfg,
		Documents: make(map[string]Document),
		Logger:    logger,
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

// For Dump
func (s *Store) Dump() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("store is nil")
	}
	return json.MarshalIndent(s, "", "  ")
}

func (s *Store) DumpToFile(filename string) error {
	data, err := s.Dump()
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
	store.Name += "Restored"
	return &store, nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}
	return NewStoreFromDump(data)
}

// For marshaling
func (s *Store) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	fields := map[string]interface{}{
		"Name":        s.Name,
		"Collections": s.Collections,
	}
	return json.Marshal(fields)
}

func (s *Store) UnmarshalJSON(data []byte) error {
	var out struct {
		Name        string                 `json:"Name"`
		Collections map[string]*Collection `json:"Collections"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	s.Name = out.Name
	s.Collections = out.Collections
	return nil
}
