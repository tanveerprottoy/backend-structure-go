package mock

import (
	"context"
	"errors"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
)

// MemoryStorage is a mock storage
// can be used to mock the repository
// for service testing
type MemoryStorage struct {
	m map[string]user.User
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{m: make(map[string]user.User)}
}

func (s *MemoryStorage) Create(ctx context.Context, e user.User, args ...any) (string, error) {
	s.m[e.ID] = e
	return e.ID, nil
}

func (s MemoryStorage) ReadMany(ctx context.Context, limit int, offset int, args ...any) ([]user.User, int64, error) {
	entities := make([]user.User, len(s.m))
	for _, v := range s.m {
		entities = append(entities, v)
	}
	return entities, int64(len(entities)), nil
}

func (s MemoryStorage) ReadOne(ctx context.Context, id string, args ...any) (user.User, error) {
	if e, ok := s.m[id]; ok {
		return e, nil
	}
	return user.User{}, errors.New("not found")
}

func (s *MemoryStorage) Update(ctx context.Context, id string, e user.User, args ...any) (int64, error) {
	if _, ok := s.m[id]; ok {
		s.m[id] = e
		return 1, nil
	}
	// not found return error
	return -1, errors.New("not found")
}

func (s *MemoryStorage) Delete(ctx context.Context, id string, args ...any) (int64, error) {
	if e, ok := s.m[id]; ok {
		e.IsArchived = true
		s.m[id] = e
		return 1, nil
	}
	// not found return error
	return -1, errors.New("not found")
}

func (s *MemoryStorage) Clear() {
	clear(s.m)
}
