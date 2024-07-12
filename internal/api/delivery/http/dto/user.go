package dto

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
)

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

func NewUserEntityDTO(id, name string, address *string, isArchived bool, createdAt, updatedAt int64) *UserEntityDTO {
	return &UserEntityDTO{
		ID:         id,
		Name:       name,
		Address:    address,
		IsArchived: isArchived,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// helper function to convert to dto entity from domain entity
func ToUserEntity(u user.User) *UserEntityDTO {
	return NewUserEntityDTO(
		u.ID,
		u.Name,
		u.Address,
		u.IsArchived,
		u.CreatedAt,
		u.UpdatedAt,
	)
}

// helper function to convert to dto entity slice from domain entity slice
func ToUserEntities(users []user.User) []UserEntityDTO {
	entityDTOs := make([]UserEntityDTO, len(users))
	for _, u := range users {
		entityDTOs = append(entityDTOs, *ToUserEntity(u))
	}
	return entityDTOs
}
