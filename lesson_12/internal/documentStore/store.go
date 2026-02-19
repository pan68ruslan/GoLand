package documentStore

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

type Store struct {
	Name        string
	Collections map[string]*Collection `json:"Collections"`
	Logger      *slog.Logger
	mu          sync.RWMutex
}

func NewStore(name string, logger *slog.Logger) *Store {
	var nm string
	if name == "" {
		nm = "DocumentStore"
	} else {
		nm = name
	}
	return &Store{
		Name:        nm,
		Collections: make(map[string]*Collection),
		Logger:      logger,
	}
}

func (s *Store) CreateCollection(name string, logger *slog.Logger) (bool, *Collection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cl, ok := s.Collections[name]; ok {
		s.Logger.Info(fmt.Sprintf("[Store]The collection '%s' already exists\n", name))
		return true, cl
	}
	col := NewCollection(name, logger)
	s.Collections[name] = &col
	s.Logger.Info(fmt.Sprintf("[Store]The collection '%s' was created\n", name))
	return true, &col
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if col, ok := s.Collections[name]; ok {
		s.Logger.Info(fmt.Sprintf("[Store]The collection '%s' was found\n", name))
		return col, true
	}
	s.Logger.Info(fmt.Sprintf("[Store]The collection '%s' was not found\n", name))
	return nil, false
}

func (s *Store) DeleteCollection(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Collections[name]; ok {
		s.Logger.Info(fmt.Sprintf("[Store]The collection '%s' has been deleted\n", name))
		delete(s.Collections, name)
		return true
	}
	s.Logger.Info(fmt.Sprintf("[Store]The collection '%s' doesn't exist\n", name))
	return false
}

func (s *Store) GetCollectionList(param ...string) string {
	result := ""
	count := 0
	if len(param) > 0 {
		if i, err := strconv.Atoi(param[0]); err == nil {
			count = i
		}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.Collections {
		if len(result) > 0 {
			result += ", "
		}
		result += c.Name
		count--
		if count == 0 {
			break
		}
	}
	return fmt.Sprintf("[%s]", result)
}

// For Dump
func (s *Store) Dump() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("store is nil")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
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
	s.mu.RLock()
	defer s.mu.RUnlock()
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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Name = out.Name
	s.Collections = out.Collections
	return nil
}
