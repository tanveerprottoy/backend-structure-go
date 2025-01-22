package user

import "errors"

type User struct {
	ID         string
	Name       string
	Address    *string
	IsArchived bool
	CreatedAt  int64
	UpdatedAt  int64
}

func MakeUser(id, name string, address *string, createdAt, updatedAt int64) User {
	u := User{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	u.setNillableFields(address)

	return u
}

func NewUser(id, name string, address *string, createdAt, updatedAt int64) *User {
	u := &User{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	u.setNillableFields(address)

	return u
}

func (u *User) setNillableFields(address *string) {
	if address != nil && *address != "" {
		u.Address = address
	}
}

func (u User) Validate() error {
	if u.Name == "" {
		return errors.New("name required")
	}

	return nil
}
