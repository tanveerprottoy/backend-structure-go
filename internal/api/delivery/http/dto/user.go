package dto

import (
	"encoding/json"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
)

type CreateUser struct {
	Name    string  `json:"name" validate:"required"`
	Address *string `json:"description" validate:"omitempty"`
}

func (u *CreateUser) ToDomainDTO() *user.CreateDTO {
	return &user.CreateDTO{
		Name:    u.Name,
		Address: u.Address,
	}
}

type UpdateUser struct {
	Name       string  `json:"name" validate:"required"`
	Address    *string `json:"description" validate:"omitempty"`
	IsArchived bool    `json:"isArchived" validate:"boolean"`
}

func (u *UpdateUser) ToDomainDTO() *user.UpdateDTO {
	return &user.UpdateDTO{
		Name:       u.Name,
		Address:    u.Address,
		IsArchived: u.IsArchived,
	}
}

// UserEntityAlias is a custom type to avoid infinite recursion in MarshalJSON
// As UserEntityAlias itself doesn't have MarshalJSON implemented,
// it doesn't infinitely recurse
type UserEntityAlias UserEntity

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

func (u UserEntity) MarshalJSON() ([]byte, error) {
	// This will not recurse infinitely
	return json.Marshal(&struct {
		UserEntityAlias
	}{
		UserEntityAlias: UserEntityAlias(u),
	})

	// or
	// return json.Marshal(UserEntityAlias(u))
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
