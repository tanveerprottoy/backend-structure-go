package dto

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
