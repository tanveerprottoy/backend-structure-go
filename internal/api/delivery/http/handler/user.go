package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/dto"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/httpext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/response"
	"github.com/tanveerprottoy/backend-structure-go/pkg/validatorext"
)

// User handles incoming requests
type User struct {
	useCase   user.UseCase
	validater validatorext.Validater
}

// NewUser initializes a new Handler
func NewUser(u user.UseCase, v validatorext.Validater) *User {
	return &User{useCase: u, validater: v}
}

// Create handles entity create post request
func (h *User) Create(w http.ResponseWriter, r *http.Request) {
	var v dto.CreateUser
	// parse the request body
	defer r.Body.Close()
	err := httpext.ParseRequestBody(r.Body, &v)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// validate the request body
	errors := h.validater.Validate(&v)
	if errors != nil {
		response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorMultiple, errors))
		return
	}

	// create the domain dto
	userDTO := user.CreateDTO{
		Name:    v.Name,
		Address: v.Address,
	}

	d, err := h.useCase.Create(r.Context(), userDTO)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entity
	p := dto.ToUserEntity(d)

	_, err = response.Respond(w, http.StatusCreated, response.MakeResponse(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *User) ReadMany(w http.ResponseWriter, r *http.Request) {
	limit := 10
	page := 1
	var err error

	limitStr := httpext.GetQueryParam(r, constant.ParamLimit)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorSingle, []error{fmt.Errorf("%s: %s", constant.InvalidQueryParam, limitStr)}))
			return
		}
	}

	pageStr := httpext.GetQueryParam(r, constant.ParamPage)
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorSingle, []error{fmt.Errorf("%s: %s", constant.InvalidQueryParam, pageStr)}))
			return
		}
	}

	var isArchived = false
	isArchivedStr := httpext.GetQueryParam(r, constant.ParamIsArchived)
	if isArchivedStr == "true" {
		isArchived = true
	}
	args := []any{isArchived}
	d, err := h.useCase.ReadMany(r.Context(), limit, page, args...)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entities
	i := dto.ToUserEntities(d)

	res := response.ReadManyResponse[dto.UserEntity]{
		Items: i,
		Limit: limit,
		Page:  page,
	}

	_, err = response.Respond(w, http.StatusOK, response.MakeResponse(&res))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *User) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorSingle, []error{errors.New(constant.MissingRequiredPathParam)}))
		return
	}

	d, err := h.useCase.ReadOne(r.Context(), id)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entity
	p := dto.ToUserEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.MakeResponse(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *User) Update(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorSingle, []error{errors.New(constant.ApiPattern)}))
		return
	}

	defer r.Body.Close()

	// parse the request body
	var v dto.UpdateUser
	err := httpext.ParseRequestBody(r.Body, &v)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// validate the request body
	errors := h.validater.Validate(&v)
	if errors != nil {
		response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorMultiple, errors))
		return
	}

	// create the domain dto
	userDTO := user.UpdateDTO{
		Name:       v.Name,
		Address:    v.Address,
		IsArchived: v.IsArchived,
	}

	d, err := h.useCase.Update(r.Context(), id, userDTO)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entity
	p := dto.ToUserEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.MakeResponse(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *User) Delete(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.MakeErrorResponse(constant.ErrorSingle, []error{errors.New(constant.MissingRequiredPathParam)}))
		return
	}

	d, err := h.useCase.Delete(r.Context(), id)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entity
	p := dto.ToUserEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.MakeResponse(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}
