package product

import (
	"context"
)

// Repository defines the data persistance logic that needs to be implemented
type Repository interface {
	Create(ctx context.Context, e *Product, args ...any) (string, error)

	ReadMany(ctx context.Context, limit, offset int, args ...any) ([]Product, error)

	ReadOne(ctx context.Context, id string, args ...any) (Product, error)

	Update(ctx context.Context, id string, e *Product, args ...any) (int64, error)

	Delete(ctx context.Context, id string, args ...any) (int64, error)
}
