package product

import (
	"context"

	"github.com/tanveerprottoy/backend-structure-go/pkg/response"
)

type UseCase interface {
	Create(ctx context.Context, d CreateDTO) (Product, error)

	ReadMany(ctx context.Context, limit, page int, args ...any) (response.ReadManyResponse[Product], error)

	ReadOne(ctx context.Context, id string) (Product, error)

	Update(ctx context.Context, id string, d UpdateDTO) (Product, error)

	Delete(ctx context.Context, id string) (Product, error)
}
