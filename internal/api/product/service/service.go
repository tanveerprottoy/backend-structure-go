package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/timeext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/util"
)

// service implements the use case of the product
// service contains the business logic as well as calls to the
// repository to perform db operations
type service struct {
	repository product.Repository
}

// NewService initializes a new Service
func NewService(r product.Repository) *service {
	return &service{repository: r}
}

// readOneInternal fetches one entity from db
func (s *service) readOneInternal(ctx context.Context, id string) (product.Product, error) {
	e, err := s.repository.ReadOne(ctx, id)
	if err != nil {
		return e, errorext.BuildCustomError(err)
	}

	return e, nil
}

// create defines the business logic for create post request
func (s *service) Create(ctx context.Context, dto *product.CreateDTO) (product.Product, error) {
	// build entity
	n := timeext.NowUnix()

	dto.CreatedAt = n
	dto.UpdatedAt = n

	l, err := s.repository.Create(ctx, dto)
	if err != nil {
		return product.Product{}, errorext.BuildCustomError(err)
	}

	return *product.NewProduct(
		l,
		dto.Name,
		dto.Description,
		dto.CreatedAt,
		dto.UpdatedAt,
	), nil
}

func (s *service) ReadMany(ctx context.Context, limit, page int, args ...any) ([]product.Product, error) {
	offset := util.CalculateOffset(limit, page)

	d, err := s.repository.ReadMany(ctx, limit, offset, args...)
	if err != nil {
		return d, errorext.BuildCustomError(err)
	}

	return d, nil
}

func (s *service) ReadOne(ctx context.Context, id string) (product.Product, error) {
	e, err := s.readOneInternal(ctx, id)
	if err != nil {
		return e, err
	}

	return e, nil
}

func (s *service) Update(ctx context.Context, id string, dto *product.UpdateDTO) (product.Product, error) {
	e, err := s.readOneInternal(ctx, id)
	if err != nil {
		return e, err
	}

	rowCount, err := s.repository.Update(ctx, id, dto)
	if err != nil {
		return e, errorext.BuildCustomError(err)
	}

	if rowCount > 0 {
		return e, nil
	}

	return e, errorext.NewCustomError(http.StatusBadRequest, errors.New(constant.GenericFailMessage))
}

func (s *service) Delete(ctx context.Context, id string) (product.Product, error) {
	e, err := s.readOneInternal(ctx, id)
	if err != nil {
		return e, err
	}

	n := timeext.NowUnix()
	rowCount, err := s.repository.Delete(ctx, id, n)
	if err != nil {
		return e, errorext.BuildCustomError(err)
	}

	if rowCount > 0 {
		e.IsArchived = true
		e.UpdatedAt = n
		return e, nil
	}

	return e, errorext.NewCustomError(http.StatusBadRequest, errors.New(constant.GenericFailMessage))
}
