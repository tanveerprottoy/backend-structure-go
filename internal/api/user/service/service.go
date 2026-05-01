package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
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
func (s *service) Create(ctx context.Context, payload user.CreateDTO) (user.User, error) {
	// build entity
	n := time.Now().Unix()

	payload.CreatedAt = n
	payload.UpdatedAt = n

	l, err := s.repository.Create(ctx, payload)
	if err != nil {
		return user.User{}, errorext.BuildCustomError(err)
	}

	return user.MakeUser(
		l,
		payload.Name,
		payload.Address,
		payload.CreatedAt,
		payload.UpdatedAt,
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

func (s *service) Update(ctx context.Context, id string, payload user.UpdateDTO) (user.User, error) {
	e, err := s.readOneInternal(ctx, id)
	if err != nil {
		return e, err
	}

	payload.UpdatedAt = time.Now().Unix()

	rowCount, err := s.repository.Update(ctx, id, payload)
	if err != nil {
		return e, errorext.BuildCustomError(err)
	}

	if rowCount > 0 {
		return user.MakeUser(
			id,
			payload.Name,
			payload.Address,
			e.CreatedAt,
			payload.UpdatedAt,
		), nil
	}

	return user.User{}, errorext.NewCustomError(http.StatusBadRequest, errors.New(constant.GenericFailMessage))
}

func (s *service) Delete(ctx context.Context, id string) (user.User, error) {
	e, err := s.readOneInternal(ctx, id)
	if err != nil {
		return e, err
	}

	n := time.Now().Unix()
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
