package dto

import "github.com/tanveerprottoy/backend-structure-go/internal/api/product"

type CreateProductDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdateProductDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	IsArchived  bool   `json:"isArchived" validate:"boolean"`
}

type ProductEntityDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsArchived  bool    `json:"isArchived"`
	CreatedAt   int64   `json:"createdAt"`
	UpdatedAt   int64   `json:"updatedAt"`
}

func NewProductEntityDTO(id, name string, description *string, isArchived bool, createdAt, updatedAt int64) *ProductEntityDTO {
	return &ProductEntityDTO{
		ID:          id,
		Name:        name,
		Description: description,
		IsArchived:  isArchived,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// helper function to convert to dto entity from domain entity
func ToProductEntity(p product.Product) *ProductEntityDTO {
	return NewProductEntityDTO(
		p.ID,
		p.Name,
		p.Description,
		p.IsArchived,
		p.CreatedAt,
		p.UpdatedAt,
	)
}

// helper function to convert to dto entity slice from domain entity slice
func ToProductEntities(products []product.Product) []ProductEntityDTO {
	entityDTOs := make([]ProductEntityDTO, len(products))
	for _, p := range products {
		entityDTOs = append(entityDTOs, *ToProductEntity(p))
	}
	return entityDTOs
}
