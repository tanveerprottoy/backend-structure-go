package dto

import "github.com/tanveerprottoy/backend-structure-go/internal/api/product"

type CreateProduct struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}

func (p *CreateProduct) ToDomainDTO() *product.CreateDTO {
	return &product.CreateDTO{
		Name:        p.Name,
		Description: p.Description,
	}
}

type UpdateProduct struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	IsArchived  bool   `json:"isArchived" validate:"boolean"`
}

func (p *UpdateProduct) ToDomainDTO() *product.UpdateDTO {
	return &product.UpdateDTO{
		Name:        p.Name,
		Description: p.Description,
		IsArchived:  p.IsArchived,
	}
}

type ProductEntity struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsArchived  bool    `json:"isArchived"`
	CreatedAt   int64   `json:"createdAt"`
	UpdatedAt   int64   `json:"updatedAt"`
}

func NewProductEntity(id, name string, description *string, isArchived bool, createdAt, updatedAt int64) *ProductEntity {
	return &ProductEntity{
		ID:          id,
		Name:        name,
		Description: description,
		IsArchived:  isArchived,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// helper function to convert to dto entity from domain entity
func ToProductEntity(p product.Product) *ProductEntity {
	return NewProductEntity(
		p.ID,
		p.Name,
		p.Description,
		p.IsArchived,
		p.CreatedAt,
		p.UpdatedAt,
	)
}

// helper function to convert to dto entity slice from domain entity slice
func ToProductEntities(products []product.Product) []ProductEntity {
	entityDTOs := make([]ProductEntity, len(products))
	for _, p := range products {
		entityDTOs = append(entityDTOs, *ToProductEntity(p))
	}

	return entityDTOs
}
