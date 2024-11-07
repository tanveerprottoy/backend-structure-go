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

func NewUser(id, name, address string, createdAt, updatedAt int64) *User {
	p := &User{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	p.SetAddress(address)
	return p
}

func (p User) Validate() error {
	if p.Name == "" {
		return errors.New("name required")
	}
	return nil
}

func (p *User) SetAddress(val string) {
	if val != "" {
		v := val
		p.Address = &v
	}
}
