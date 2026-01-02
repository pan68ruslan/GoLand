package documentstore

import (
	"errors"
	"fmt"
)

var (
	ErrCannotCreateCollection  = errors.New("cannot create the collection")
	ErrCannotFindCollection    = errors.New("cannot find the collection")
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrWrongCollectionConfig   = errors.New("collection config is absent or wrong")
)

type Store struct {
	Name        string
	Collections map[string]*Collection
}

func NewStore(storeName string) *Store {
	if storeName == "" {
		storeName = "Noname_Store"
	}
	return &Store{
		Name:        storeName,
		Collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, config *CollectionConfig) (*Collection, error) {
	if col, ok := s.Collections[name]; ok {
		return col, fmt.Errorf("%w: name '%s'", ErrCollectionAlreadyExists, name)
	}
	if config == nil || len(config.PrimaryKey) == 0 {
		return nil, fmt.Errorf("%w: name '%s'", ErrWrongCollectionConfig, name)
	}
	if col, e := NewCollection(name, config); e == nil {
		s.Collections[name] = col
		return col, nil
	} else {
		return nil, fmt.Errorf("%w: '%v'", ErrCannotCreateCollection, e)
	}
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	if col, ok := s.Collections[name]; ok {
		return col, nil
	}
	return nil, fmt.Errorf("%w: collection '%s'", ErrCannotFindCollection, name)
}

func (s *Store) DeleteCollection(name string) error {
	if _, ok := s.Collections[name]; ok {
		delete(s.Collections, name)
		return nil
	}
	return fmt.Errorf("%w: collection '%s'", ErrCannotFindCollection, name)
}
