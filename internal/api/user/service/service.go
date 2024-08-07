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
func (s *service) Create(ctx context.Context, d user.CreateDTO) (user.User, error) {
	// build entity
	n := timeext.NowUnix()
	e := *user.NewUser("", d.Name, d.Address, n, n)
	// check if product is valid
	err := e.Validate()
	if err != nil {
		return e, errorext.NewCustomError(http.StatusBadRequest, err)
	}

	l, err := s.repository.Create(ctx, e)
	if err != nil {
		return e, errorext.BuildCustomError(err)
	}

	e.ID = l
	return e, nil
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

func (s *service) Update(ctx context.Context, id string, d user.UpdateDTO) (user.User, error) {
	e, err := s.readOneInternal(ctx, id)
	if err != nil {
		return e, err
	}

	e.Name = d.Name
	e.UpdatedAt = timeext.NowUnix()
	// set description
	e.SetAddress(d.Address)
	rowCount, err := s.repository.Update(ctx, id, e)
	if err != nil {
		return e, errorext.BuildCustomError(err)
	}

	if rowCount > 0 {
		return e, nil
	}

	return e, errorext.NewCustomError(http.StatusBadRequest, errors.New(constant.GenericFailMessage))
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
