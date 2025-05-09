package mock

import (
	"context"
	"errors"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
)

// MemoryStorage is a mock storage
// can be used to mock the repository
// for service testing
type MemoryStorage struct {
	m map[string]*product.Product
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{m: make(map[string]*product.Product)}
}

func (s *MemoryStorage) Create(ctx context.Context, dto *product.CreateDTO, args ...any) (string, error) {
	s.m[dto.Name] = product.NewProduct(dto.Name, dto.Name, nil, 0, 0)

	return dto.Name, nil
}

func (s MemoryStorage) ReadMany(ctx context.Context, limit int, offset int, args ...any) ([]product.Product, error) {
	entities := make([]product.Product, len(s.m))

	for _, v := range s.m {
		entities = append(entities, *v)
	}

	return entities, nil
}

func (s MemoryStorage) ReadOne(ctx context.Context, id string, args ...any) (product.Product, error) {
	if e, ok := s.m[id]; ok {
		return *e, nil
	}

	return product.Product{}, errors.New("not found")
}

func (s *MemoryStorage) Update(ctx context.Context, id string, dto *product.UpdateDTO, args ...any) (int64, error) {
	if e, ok := s.m[id]; ok {
		e.Name = dto.Name
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
