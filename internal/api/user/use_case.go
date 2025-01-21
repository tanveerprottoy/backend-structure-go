package user

import (
	"context"
)

type UseCase interface {
	Create(ctx context.Context, d *CreateDTO) (User, error)

	ReadMany(ctx context.Context, limit, page int, args ...any) ([]User, error)

	ReadOne(ctx context.Context, id string) (User, error)

	Update(ctx context.Context, id string, d *UpdateDTO) (User, error)

	Delete(ctx context.Context, id string) (User, error)
}
