package user

import (
	"context"

	"github.com/tanveerprottoy/backend-structure-go/pkg/response"
)

type UseCase interface {
	Create(ctx context.Context, d CreateDTO) (User, error)

	ReadMany(ctx context.Context, limit, page int, args ...any) (response.ReadManyResponse[User], error)

	ReadOne(ctx context.Context, id string) (User, error)

	Update(ctx context.Context, id string, d UpdateDTO) (User, error)

	Delete(ctx context.Context, id string) (User, error)
}
