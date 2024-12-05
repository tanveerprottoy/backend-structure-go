package dto

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
)

type CreateUser struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"description" validate:"omitempty"`
}

type UpdateUser struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"description" validate:"omitempty"`
	IsArchived bool   `json:"isArchived" validate:"boolean"`
}

type UserEntity struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Address    *string `json:"description"`
	IsArchived bool    `json:"isArchived"`
	CreatedAt  int64   `json:"createdAt"`
	UpdatedAt  int64   `json:"updatedAt"`
}

func NewUserEntity(id, name string, address *string, isArchived bool, createdAt, updatedAt int64) *UserEntity {
	return &UserEntity{
		ID:         id,
		Name:       name,
		Address:    address,
		IsArchived: isArchived,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// helper function to convert to dto entity from domain entity
func ToUserEntity(u user.User) *UserEntity {
	return NewUserEntity(
		u.ID,
		u.Name,
		u.Address,
		u.IsArchived,
		u.CreatedAt,
		u.UpdatedAt,
	)
}

// helper function to convert to dto entity slice from domain entity slice
func ToUserEntities(users []user.User) []UserEntity {
	entityDTOs := make([]UserEntity, len(users))
	for _, u := range users {
		entityDTOs = append(entityDTOs, *ToUserEntity(u))
	}
	return entityDTOs
}
