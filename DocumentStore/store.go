package DocumentStore

import "fmt"

type Store struct {
	collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	if _, ok := s.collections[name]; ok {
		fmt.Printf("[Store]The collection '%s' already exists\n", name)
		return false, nil
	}
	col := &Collection{
		Name:      name,
		Cfg:       cfg,
		documents: make(map[string]Document),
	}
	s.collections[name] = col
	fmt.Printf("[Store]The collection '%s' created\n", name)
	return true, col
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	if col, ok := s.collections[name]; ok {
		fmt.Printf("[Store]The collection '%s' was found\n", name)
		return col, true
	}
	fmt.Printf("[Store]The collection '%s' was not found\n", name)
	return nil, false
}

func (s *Store) DeleteCollection(name string) bool {
	if _, ok := s.collections[name]; ok {
		fmt.Printf("[Store]The collection '%s' has been deleted\n", name)
		delete(s.collections, name)
		return true
	}
	fmt.Printf("[Store]The collection '%s' doesn't exist\n", name)
	return false
}
