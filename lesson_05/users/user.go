package users

import (
	"errors"
	"fmt"

	store "github.com/pan68ruslan/GoLand/lesson_05/documentstore"
)

var (
	ErrCantCreateService     = errors.New("cannot create the user service")
	ErrUserCreating          = errors.New("cannot add the user")
	ErrUnmarshallingDocument = errors.New("cannot unmarshal the document")
	ErrUserRemoving          = errors.New("cannot remove the user")
	ErrUserNotFound          = errors.New("user not found")
	ErrBrokenUser            = errors.New("user is broken or unknown")
	ErrListingUsers          = errors.New("users listing error: ")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	collection *store.Collection
}

func NewService(name string, s *store.Store) (*Service, error) {
	if col, e := s.CreateCollection(name, &store.CollectionConfig{PrimaryKey: "id"}); e == nil {
		return &Service{collection: col}, nil
	} else {
		return nil, fmt.Errorf("%w: error: %v", ErrCantCreateService, e)
	}
}

func (s *Service) CreateUser(u User) (*User, error) {
	if doc, e := store.MarshalStructureToDocument(u); e == nil {
		s.collection.Documents[u.ID] = doc
		return &u, nil
	} else {
		return nil, fmt.Errorf("%w, error %v: userId=%s, userName=%s", ErrUserCreating, e, u.ID, u.Name)
	}
}

func (s *Service) GetUser(id string) (User, error) {
	doc, e := s.collection.Get(id)
	if e != nil {
		return User{}, fmt.Errorf("%w, error %v: userId=%s", ErrUserNotFound, e, id)
	}
	var u User
	if e = store.UnmarshalDocumentToStructure(&doc, &u); e != nil {
		return User{}, fmt.Errorf("%w, error %v: userId=%s", ErrUnmarshallingDocument, e, id)
	}
	return u, nil
}

func (s *Service) DeleteUser(id string) error {
	if e := s.collection.Delete(id); e != nil {
		return fmt.Errorf("%w, error %v: userId=%s", ErrUserRemoving, e, id)
	}
	return nil
}

func (s *Service) ListUsers() ([]User, error) {
	var errs []error
	var users []User
	errs = append(errs, ErrListingUsers)
	for _, doc := range s.collection.Documents {
		var u User
		e := store.UnmarshalDocumentToStructure(&doc, &u)
		if e == nil {
			users = append(users, u)
		} else {
			errs = append(errs, fmt.Errorf("%w, error=%v", ErrBrokenUser, e.Error()))
		}
	}
	if len(errs) > 1 {
		return nil, errors.Join(errs...)
	}
	return users, nil
}
