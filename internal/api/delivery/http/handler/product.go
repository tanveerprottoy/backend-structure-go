package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/dto"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/httpext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/response"
	"github.com/tanveerprottoy/backend-structure-go/pkg/validatorext"
)

// Product handles incoming requests
// it validates request body
// converts handler dto to domain dto
// converts domain entity to handler response payload
type Product struct {
	useCase   product.UseCase
	validater validatorext.Validater
}

// NewProduct initializes a new Handler
func NewProduct(u product.UseCase, v validatorext.Validater) *Product {
	return &Product{useCase: u, validater: v}
}

// Create handles entity create post request
func (h *Product) Create(w http.ResponseWriter, r *http.Request) {
	var v dto.CreateProduct
	err := httpext.ParseRequestBody(r.Body, &v)
	if err != nil {
		err = errorext.ParseJSONError(err)
		response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// validate the request body
	errs := h.validater.Validate(&v)
	if errs != nil {
		response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorMultiple, errs))
		return
	}

	// create the domain dto
	productDTO := v.ToDomainDTO()

	d, err := h.useCase.Create(r.Context(), productDTO)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), response.NewErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// convert to dto entity
	p := dto.ToProductEntity(d)

	_, err = response.Respond(w, http.StatusCreated, response.NewResponse(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *Product) ReadMany(w http.ResponseWriter, r *http.Request) {
	// temp
	response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorSingle, []error{errors.New("bad request")}))
	return

	limit := 10
	page := 1
	var err error

	limitStr := httpext.GetQueryParam(r, constant.ParamLimit)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorSingle, []error{fmt.Errorf("%s: %s", constant.InvalidQueryParam, limitStr)}))
			return
		}
	}

	pageStr := httpext.GetQueryParam(r, constant.ParamPage)
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorSingle, []error{fmt.Errorf("%s: %s", constant.InvalidQueryParam, pageStr)}))
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
		response.RespondError(w, err.Code(), response.NewErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// convert to dto entities
	i := dto.ToProductEntities(d)

	res := response.ReadManyResponse[dto.ProductEntity]{
		Items: i,
		Limit: limit,
		Page:  page,
	}

	_, err = response.Respond(w, http.StatusOK, response.NewResponse(&res))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *Product) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorSingle, []error{errors.New(constant.MissingRequiredPathParam)}))
		return
	}

	d, err := h.useCase.ReadOne(r.Context(), id)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), response.NewErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// convert to dto entity
	p := dto.ToProductEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.NewResponse(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *Product) Update(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorSingle, []error{errors.New(constant.ApiPattern)}))
		return
	}

	// parse the request body
	var v dto.UpdateProduct
	err := httpext.ParseRequestBody(r.Body, &v)
	if err != nil {
		err = errorext.ParseJSONError(err)
		response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// validate the request body
	errs := h.validater.Validate(&v)
	if errs != nil {
		response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorMultiple, errs))
		return
	}

	// create the domain dto
	productDTO := v.ToDomainDTO()

	d, err := h.useCase.Update(r.Context(), id, productDTO)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), response.NewErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// convert to dto entity
	p := dto.ToProductEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.NewResponse(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *Product) Delete(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.NewErrorResponse(constant.ErrorSingle, []error{errors.New(constant.MissingRequiredPathParam)}))
		return
	}

	d, err := h.useCase.Delete(r.Context(), id)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), response.NewErrorResponse(constant.ErrorSingle, []error{err}))
		return
	}

	// convert to dto entity
	p := dto.ToProductEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.NewResponse(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}
