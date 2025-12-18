package users

import (
	"errors"
	"fmt"
)

var (
	ErrUserIdIsEmpty     = errors.New("userId cannot be empty")
	ErrNameIsEmpty       = errors.New("name cannot be empty")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Collection struct {
	Title string
	Users map[string]*User
}

type Service struct {
	coll *Collection
}

func NewService(coll *Collection) *Service {
	if coll.Users == nil {
		coll.Users = make(map[string]*User)
	}
	return &Service{coll: coll}
}

func (s *Service) CreateUser(id, name string) (*User, error) {
	var errs []error
	if id == "" {
		errs = append(errs, ErrUserIdIsEmpty)
	}
	if name == "" {
		errs = append(errs, ErrNameIsEmpty)
	}
	if s.coll.Users == nil {
		s.coll.Users = make(map[string]*User)
	}
	if _, exists := s.coll.Users[id]; exists {
		errs = append(errs, fmt.Errorf("%w: id=%s", ErrUserAlreadyExists, id))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	u := &User{
		ID:   id,
		Name: name,
	}
	s.coll.Users[id] = u
	return u, nil
}

func (s *Service) ListUsers() ([]User, error) {
	if s.coll.Users == nil {
		return []User{}, nil
	}
	result := make([]User, 0, len(s.coll.Users))
	for _, u := range s.coll.Users {
		result = append(result, *u)
	}
	return result, nil
}

func (s *Service) GetUser(userId string) (*User, error) {
	var errs []error
	if userId == "" {
		errs = append(errs, ErrUserIdIsEmpty)
	}
	u, ok := s.coll.Users[userId]
	if !ok {
		errs = append(errs, fmt.Errorf("%w: id=%s", ErrUserNotFound, userId))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return u, nil
}

func (s *Service) DeleteUser(userId string) error {
	if _, e := s.GetUser(userId); e != nil {
		delete(s.coll.Users, userId)
	}
	return nil
}
