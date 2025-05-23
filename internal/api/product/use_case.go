package product

import (
	"context"
)

type UseCase interface {
	Create(ctx context.Context, dto *CreateDTO) (Product, error)

	ReadMany(ctx context.Context, limit, page int, args ...any) ([]Product, error)

	ReadOne(ctx context.Context, id string) (Product, error)

	Update(ctx context.Context, id string, dto *UpdateDTO) (Product, error)

	Delete(ctx context.Context, id string) (Product, error)
}
