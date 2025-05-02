package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/timeext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/util"
)

// service implements the use case for user
// service contains the business logic as well as calls to the
// repository to perform db operations
type service struct {
	repository user.Repository
}

// NewService initializes a new Service
func NewService(r user.Repository) *service {
	return &service{repository: r}
}

// readOneInternal fetches one entity from db
func (s *service) readOneInternal(ctx context.Context, id string) (user.User, error) {
	e, err := s.repository.ReadOne(ctx, id)
	if err != nil {
		return e, errorext.BuildCustomError(err)
	}

	return e, nil
}

// create defines the business logic for create post request
func (s *service) Create(ctx context.Context, dto *user.CreateDTO) (user.User, error) {
	// build entity
	n := timeext.NowUnix()

	dto.CreatedAt = n
	dto.UpdatedAt = n

	l, err := s.repository.Create(ctx, dto)
	if err != nil {
		return user.User{}, errorext.BuildCustomError(err)
	}

	return user.MakeUser(
		l,
		dto.Name,
		dto.Address,
		dto.CreatedAt,
		dto.UpdatedAt,
	), nil
}

func (s *service) ReadMany(ctx context.Context, limit, page int, args ...any) ([]user.User, error) {
	offset := util.CalculateOffset(limit, page)

	d, err := s.repository.ReadMany(ctx, limit, offset, args...)
	if err != nil {
		return d, errorext.BuildCustomError(err)
	}

	return d, nil
}

func (s *service) ReadOne(ctx context.Context, id string) (user.User, error) {
	e, err := s.readOneInternal(ctx, id)
	if err != nil {
		return e, err
	}

	return e, nil
}

func (s *service) Update(ctx context.Context, id string, dto *user.UpdateDTO) (user.User, error) {
	e, err := s.readOneInternal(ctx, id)
	if err != nil {
		return e, err
	}

	dto.UpdatedAt = timeext.NowUnix()

	rowCount, err := s.repository.Update(ctx, id, dto)
	if err != nil {
		return e, errorext.BuildCustomError(err)
	}

	if rowCount > 0 {
		return user.MakeUser(
			id,
			dto.Name,
			dto.Address,
			e.CreatedAt,
			dto.UpdatedAt,
		), nil
	}

	return user.User{}, errorext.NewCustomError(http.StatusBadRequest, errors.New(constant.GenericFailMessage))
}

func (s *service) Delete(ctx context.Context, id string) (user.User, error) {
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
