package user

import (
	"context"
)

// Repository defines the data persistance logic that needs to be implemented
type Repository interface {
	Create(ctx context.Context, e *User, args ...any) (string, error)

	ReadMany(ctx context.Context, limit, offset int, args ...any) ([]User, error)

	ReadOne(ctx context.Context, id string, args ...any) (User, error)

	Update(ctx context.Context, id string, e *User, args ...any) (int64, error)

	Delete(ctx context.Context, id string, args ...any) (int64, error)
}
