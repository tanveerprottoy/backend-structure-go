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
	// parse the request body
	defer r.Body.Close()
	err := httpext.ParseRequestBody(r.Body, &v)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, response.BuildError(err))
		return
	}

	// validate the request body
	validationErrs := h.validater.Validate(&v)
	if validationErrs != nil {
		response.RespondError(w, http.StatusBadRequest, response.BuildError(validationErrs))
		return
	}

	// create the domain dto
	productDTO := product.CreateDTO{
		Name:        v.Name,
		Description: v.Description,
	}

	d, err := h.useCase.Create(r.Context(), productDTO)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entity
	p := dto.ToProductEntity(d)

	_, err = response.Respond(w, http.StatusCreated, response.BuildData(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *Product) ReadMany(w http.ResponseWriter, r *http.Request) {
	limit := 10
	page := 1
	var err error

	limitStr := httpext.GetQueryParam(r, constant.ParamLimit)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			response.RespondError(w, http.StatusBadRequest, response.BuildError(fmt.Errorf(constant.InvalidQueryParam+": %s", limitStr)))
			return
		}
	}

	pageStr := httpext.GetQueryParam(r, constant.ParamPage)
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			response.RespondError(w, http.StatusBadRequest, response.BuildError(fmt.Errorf(constant.InvalidQueryParam+": %s", pageStr)))
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
	i := dto.ToProductEntities(d)

	res := response.ReadManyResponse[dto.ProductEntity]{
		Items: i,
		Limit: limit,
		Page:  page,
	}

	_, err = response.Respond(w, http.StatusOK, response.BuildData(&res))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *Product) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.BuildError(errors.New(constant.MissingRequiredPathParam)))
		return
	}

	d, err := h.useCase.ReadOne(r.Context(), id)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entity
	p := dto.ToProductEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.BuildData(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *Product) Update(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.BuildError(errors.New(constant.ApiPattern)))
		return
	}

	defer r.Body.Close()

	// parse the request body
	var v dto.UpdateProduct
	err := httpext.ParseRequestBody(r.Body, &v)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, response.BuildError(err))
		return
	}

	// validate the request body
	validationErrs := h.validater.Validate(&v)
	if validationErrs != nil {
		response.RespondError(w, http.StatusBadRequest, response.BuildError(validationErrs))
		return
	}

	// create the domain dto
	productDTO := product.UpdateDTO{
		Name:        v.Name,
		Description: v.Description,
		IsArchived:  v.IsArchived,
	}

	d, err := h.useCase.Update(r.Context(), id, productDTO)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entity
	p := dto.ToProductEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.BuildData(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}

func (h *Product) Delete(w http.ResponseWriter, r *http.Request) {
	id := httpext.GetURLParam(r, constant.ParamId)
	if id == "" {
		response.RespondError(w, http.StatusBadRequest, response.BuildError(errors.New(constant.MissingRequiredPathParam)))
		return
	}

	d, err := h.useCase.Delete(r.Context(), id)
	if err != nil {
		err := errorext.ParseCustomError(err)
		response.RespondError(w, err.Code(), err.Err())
		return
	}

	// convert to dto entity
	p := dto.ToProductEntity(d)

	_, err = response.Respond(w, http.StatusOK, response.BuildData(p))
	if err != nil {
		log.Printf("response.Respond returned error: %v", err)
	}
}
