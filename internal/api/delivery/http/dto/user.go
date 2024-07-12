package dto

type CreateUserDTO struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"description" validate:"omitempty"`
}

type UpdateUserDTO struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"description" validate:"omitempty"`
	IsArchived bool   `json:"isArchived" validate:"boolean"`
}

type UserEntityDTO struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Address    *string `json:"description"`
	IsArchived bool    `json:"isArchived"`
	CreatedAt  int64   `json:"createdAt"`
	UpdatedAt  int64   `json:"updatedAt"`
}
